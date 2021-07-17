package bucket

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/minio/minio-go"
)

const DOAccessKey = "xxx"
const DOSecretAccessKey = "xxx"
const DOEndpoint = "fra1.digitaloceanspaces.com"
const bucketName = "grbpwr"
const objectName = "test.png"
const filePath = "./test.png"
const contentType = "image/png"

const b64Image = ""

func imageToB64(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	return base64Encoding, nil
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func TestClient(t *testing.T) {
	client, err := minio.New(DOEndpoint, DOAccessKey, DOSecretAccessKey, true)
	if err != nil {
		log.Fatal(err)
	}

	spaces, err := client.ListBuckets()
	if err != nil {
		log.Fatal("list err ", err)
	}

	// err = client.MakeBucket("grbpwr-com", "fra-1")
	// if err != nil {
	// 	log.Fatal("MakeBucket err ", err)
	// }

	for _, space := range spaces {
		fmt.Println(space.Name)
	}

	i, err := imageToB64(filePath)
	if err != nil {
		log.Fatal("imageToB64 err ", err)
	}

	r, ft, err := B64ToImage(i)

	fmt.Println(ft.Extension)
	fmt.Println(ft.MIMEType)
	fmt.Printf("client %+v ", client)

	_, err = client.PutObject(bucketName, objectName, r, r.Size(), minio.PutObjectOptions{ContentType: contentType})

}

func TestUploadImage(t *testing.T) {
	b := &Bucket{}
	err := env.Parse(b)
	if err != nil {
		log.Fatal("Parse err ", err)
	}

	err = b.GetBucket()
	if err != nil {
		log.Fatal("GetBucket err ", err)
	}

	i, err := imageToB64(filePath)
	if err != nil {
		log.Fatal("imageToB64 err ", err)
	}

	b.DOAccessKey = DOAccessKey
	b.DOSecretAccessKey = DOSecretAccessKey
	b.DOEndpoint = DOEndpoint

	fp, err := b.UploadImage(i)
	fmt.Println("--- ", fp)
	fmt.Println("--- ", err)

}
