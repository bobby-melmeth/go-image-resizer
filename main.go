package main

import (
	"bytes"
	"context"
	"image/jpeg"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

func Handler(ctx context.Context, event events.S3Event) {

	awsZone := os.Getenv("AWS_ZONE")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsZone),
	}))

	s3Client := s3.New(sess)

	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		obj, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})

		if err != nil {
			log.Printf("Error getting object %s from bucket %s: %v", key, bucket, err)
			return
		}

		img, err := jpeg.Decode(obj.Body)
		if err != nil {
			log.Printf("Error decoding image %s: %v", key, err)
			return
		}

		newImg := resize.Resize(200, 0, img, resize.Lanczos3)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, newImg, nil)
		if err != nil {
			log.Printf("Error encoding resized image %s: %v", key, err)
			return
		}

		newKey := strings.Replace(key, "current/", "resized/", 1)
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: &bucket,
			Key:    &newKey,
			Body:   bytes.NewReader(buf.Bytes()),
		})
		if err != nil {
			log.Printf("Error uploading resized image %s to bucket %s: %v", newKey, bucket, err)
			return
		}

		log.Printf("Image %s resized and saved to %s", key, newKey)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.SetFlags(log.LstdFlags | log.Llongfile)
	lambda.Start(Handler)
}
