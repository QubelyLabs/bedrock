package db

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"strings"

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

func Upload(bucketName string, objectName string, objectContent io.Reader, contentType string, encoding string) error {
	s3 := S3()
	ctx := context.TODO()

	exists, err := s3.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err := s3.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			return err
		}
	}

	buf := &bytes.Buffer{}
	size, err := io.Copy(buf, objectContent)
	if err != nil {
		return err
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if encoding == "" {
		encoding = "base64"
	}

	_, err = s3.PutObject(ctx, bucketName, objectName, buf, size, minio.PutObjectOptions{
		ContentType:     contentType,
		ContentEncoding: encoding,
	})

	if err != nil {
		return err
	}

	return nil
}

func Download(bucketName, objectName string) (string, string, error) {
	object, err := s3.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", "", err
	}
	defer object.Close()

	// Read object metadata to get Content-Type
	info, err := object.Stat()
	if err != nil {
		return "", "", err
	}
	contentType := info.ContentType

	// Read object data and encode it to base64
	var buf strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := io.Copy(encoder, object); err != nil {
		return "", "", err
	}
	encoder.Close() // Ensure all data is flushed

	return buf.String(), contentType, nil
}
