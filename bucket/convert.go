package bucket

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

func PNGFromB64(b64Image []byte) (image.Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b64Image))
	i, err := png.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("PNGFromB64:image.Decode: [%v]", err.Error())
	}
	return i, nil
}

func JPGFromB64(b64Image []byte) (image.Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b64Image))
	i, err := jpeg.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("JPGFromB64:image.Decode")
	}
	return i, nil
}

func EncodeJPG(w io.Writer, img image.Image) error {
	var err error

	newImage := resize.Resize(1000, 1000, img, resize.Lanczos3)

	var rgba *image.RGBA
	if nrgba, ok := newImage.(*image.NRGBA); ok {
		if nrgba.Opaque() {
			rgba = &image.RGBA{
				Pix:    nrgba.Pix,
				Stride: nrgba.Stride,
				Rect:   nrgba.Rect,
			}
		}
	}

	if rgba != nil {
		err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: 60})
	} else {
		err = jpeg.Encode(w, newImage, &jpeg.Options{Quality: 60})
	}

	return err
}
