package db

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var s3 *minio.Client

func InitS3(endpoint, port, accessKeyID, secretAccessKey, useSSL string) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL == "true",
	})

	if err != nil {
		log.Fatalln("minio init ", err)
	}

	s3 = client
}

func S3() *minio.Client {
	return s3
}

func Upload(bucketName string, objectName string, objectContent io.Reader, mimeType string, encoding string) error {
	s3 := S3()
	ctx := context.TODO()
	err := s3.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		exists, errBucketExists := s3.BucketExists(ctx, bucketName)
		if errBucketExists != nil || !exists {
			return err
		}
	}

	buf := &bytes.Buffer{}
	size, err := io.Copy(buf, objectContent)
	if err != nil {
		return err
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if encoding == "" {
		encoding = "base64"
	}

	_, err = s3.PutObject(ctx, bucketName, objectName, objectContent, size, minio.PutObjectOptions{
		ContentType:     mimeType,
		ContentEncoding: encoding,
	})

	if err != nil {
		return err
	}

	return nil
}

func Download(bucketName, objectName string) (string, error) {
	object, err := s3.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer object.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, object); err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return base64String, nil
}
