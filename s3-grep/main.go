package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func listBucketObjects(bucketName string, client *s3.Client) {
	input := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}

	var fileKeys []string

	paginator := s3.NewListObjectsV2Paginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to list objects in bucket %s: %v", bucketName, err)
		}

		for _, obj := range page.Contents {
			fileKeys = append(fileKeys, *obj.Key)
		}
	}

	fmt.Printf("Objects in bucket '%s':\n", bucketName)
	for _, key := range fileKeys {
		fmt.Println(key)
	}
}

func main() {
	// Define a command-line flag for the S3 bucket name
	bucketName := flag.String("bucket", "", "Name of the S3 bucket")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	// If a bucket name is provided as a flag, list objects in that bucket
	if *bucketName != "" {
		listBucketObjects(*bucketName, client)
	} else {
		// If no bucket name is provided, list objects in all S3 buckets in the account
		fmt.Println("Listing objects in all S3 buckets in the account:")
		bucketsInput := &s3.ListBucketsInput{}
		bucketsOutput, err := client.ListBuckets(context.TODO(), bucketsInput)
		if err != nil {
			log.Fatalf("failed to list buckets, %v", err)
		}

		for _, bucket := range bucketsOutput.Buckets {
			listBucketObjects(*bucket.Name, client)
			fmt.Println()
		}
	}
}