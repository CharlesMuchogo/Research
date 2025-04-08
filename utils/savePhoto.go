package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SavePhoto(file *multipart.FileHeader, userID string) (string, error) {

	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_BUCKET")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return "", err
	}
	svc := s3.New(sess)

	fileExt := filepath.Ext(file.Filename)
	currentTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
	key := fmt.Sprintf("%d_%s%s", currentTimestamp, userID, fileExt)

	src, err := file.Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer src.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(mime.TypeByExtension(fileExt)), // Optional MIME type detection
	})
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key), nil

}
