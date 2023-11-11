package bucket

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/jekabolt/grbpwr-manager/internal/entity"
	pb_common "github.com/jekabolt/grbpwr-manager/proto/gen/common"
	"github.com/minio/minio-go/v7"
	"golang.org/x/image/draw"
)

type B64Image struct {
	content     []byte
	contentType string
}

// upload image to bucket return url
func (b *Bucket) uploadImageToBucket(ctx context.Context, img io.Reader, folder, imageName string, contentType ContentType) (string, error) {
	ext, err := fileExtensionFromContentType(contentType)
	if err != nil {
		return "", fmt.Errorf("can't get file extension")
	}
	fp := b.constructFullPath(folder, imageName, ext)

	data, err := io.ReadAll(img)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(data)
	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	cacheControl := "max-age=31536000"

	_, err = b.Client.PutObject(ctx, b.Config.S3BucketName, fp, r,
		int64(r.Len()), minio.PutObjectOptions{
			ContentType:  contentType.String(),
			CacheControl: cacheControl,
			UserMetadata: userMetaData,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error putting object: %v", err)
	}

	return b.getCDNURL(fp), nil
}

// getB64ImageFromString extracts the content type and the byte content from a raw base64 image string.
// The expected format of the raw base64 string is "data:[<mediatype>];base64,[<base64-data>]".
func getB64ImageFromString(rawB64Image string) (*B64Image, error) {
	const base64Prefix = ";base64,"
	parts := strings.Split(rawB64Image, base64Prefix)

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid base64 image format: expected 'data:[mediatype];base64,[data]'")
	}

	return &B64Image{
		contentType: parts[0],
		content:     []byte(parts[1]),
	}, nil
}

func (b64Img *B64Image) b64ToImage() (image.Image, error) {
	switch b64Img.contentType {
	case "data:image/jpeg":
		return decodeImageFromB64(b64Img.content, contentTypeJPEG)
	case "data:image/png":
		return decodeImageFromB64(b64Img.content, contentTypePNG)
	default:
		return nil, fmt.Errorf("b64ToImage: File type is not supported [%s]", b64Img.contentType)
	}
}
func imageFromString(rawB64Image string) (image.Image, error) {
	b64Img, err := getB64ImageFromString(rawB64Image)
	if err != nil {
		return nil, err
	}
	return b64Img.b64ToImage()
}

// upload single image with defined quality and	prefix to bucket
func (b *Bucket) uploadSingleImage(ctx context.Context, img image.Image, quality int, folder, imageName string) (string, error) {
	var buf bytes.Buffer

	// Encode the image to JPEG format with given quality.
	if err := encodeJPG(&buf, img, quality); err != nil {
		return "", fmt.Errorf("failed to encode JPG: %v", err)
	}

	// Upload the JPEG data to S3 bucket.
	url, err := b.uploadImageToBucket(ctx, &buf, folder, imageName, contentTypeJPEG)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to bucket: %v", err)
	}

	return url, nil
}

// compose internal image object (with FullSize & Compressed formats) and upload it to S3
func (b *Bucket) uploadImageObj(ctx context.Context, img image.Image, folder, imageName string) (*pb_common.Media, error) {
	imgObj := &pb_common.MediaInsert{}

	fullSizeName := fmt.Sprintf("%s-%s", imageName, "og")
	compressedName := fmt.Sprintf("%s-%s", imageName, "compressed")
	thumbnailName := fmt.Sprintf("%s-%s", imageName, "thumb")
	var err error

	// Upload full size image
	imgObj.FullSize, err = b.uploadSingleImage(ctx, img, 100, folder, fullSizeName)
	if err != nil {
		return nil, fmt.Errorf("failed to upload full-size image: %v", err)
	}

	// Upload compressed image
	imgObj.Compressed, err = b.uploadSingleImage(ctx, img, 60, folder, compressedName)
	if err != nil {
		return nil, fmt.Errorf("failed to upload compressed image: %v", err)
	}

	imgObj.Thumbnail, err = b.uploadSingleImage(ctx, resizeImage(img, 1080), 90, folder, thumbnailName)
	if err != nil {
		return nil, fmt.Errorf("failed to upload compressed image: %v", err)
	}

	mediaId, err := b.ms.AddMedia(ctx, &entity.MediaInsert{
		FullSize:   imgObj.FullSize,
		Thumbnail:  imgObj.Thumbnail,
		Compressed: imgObj.Compressed,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add media to db: %v", err)
	}

	return &pb_common.Media{
		Id:    int32(mediaId),
		Media: imgObj,
	}, nil
}

// resizeImage checks the height of the given image. If it's greater than minWidth in px,
// it resizes the image to have a height of 'minWidth' while maintaining the aspect ratio.
func resizeImage(img image.Image, minWidth int) image.Image {
	bounds := img.Bounds()

	// Check if the height is greater than 1080px
	if bounds.Dy() > minWidth {
		// Calculate new width to maintain aspect ratio
		newWidth := minWidth * bounds.Dx() / bounds.Dy()

		// Create a new image with the desired dimensions
		newImg := image.NewRGBA(image.Rect(0, 0, newWidth, minWidth))

		// Resize the image using high-quality resampling
		draw.CatmullRom.Scale(newImg, newImg.Bounds(), img, bounds, draw.Over, nil)

		return newImg
	}

	// Return the original image if no resizing is needed
	return img
}
