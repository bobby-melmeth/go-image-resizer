package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	// Create a sample S3 event
	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: "bobbysimagebucket"},
					Object: events.S3Object{Key: "current/CarlyleHat.png"},
				},
			},
		},
	}

	// Call your Handler function with the sample event
	Handler(context.TODO(), s3Event)

	// Add assertions to check if the image was resized and saved correctly
	// You can check if the resized image was uploaded to the expected path in the bucket.

	// Add more test cases as needed to cover different scenarios.

	// For example, you can check for error conditions and handle them accordingly.
}
