package main

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()

	// Initialize MinIO client
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false, // Set to true if using HTTPS
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	log.Println("Successfully connected to MinIO")

	// 1. Create a bucket
	bucketName := "my-test-bucket"
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check if bucket already exists
		exists, err := minioClient.BucketExists(ctx, bucketName)
		if err == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	} else {
		log.Printf("Successfully created bucket: %s\n", bucketName)
	}

	// 2. Upload a file
	objectName := "test-file.txt"
	filePath := "./test-file.txt"

	// Upload the file
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	// 3. List objects in bucket
	log.Println("\nListing objects in bucket:")
	objectsCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for object := range objectsCh {
		if object.Err != nil {
			log.Printf("Error listing objects: %v", object.Err)
			continue
		}
		log.Printf(" - %s (size: %d, last modified: %s)\n",
			object.Key, object.Size, object.LastModified.Format("2006-01-02 15:04:05"))
	}

	// 4. Download the file
	downloadPath := "./downloaded-file.txt"
	err = minioClient.FGetObject(ctx, bucketName, objectName, downloadPath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
	log.Printf("Successfully downloaded file to: %s\n", downloadPath)

	// 5. Check if file exists and read its content
	content, err := os.ReadFile(downloadPath)
	if err != nil {
		log.Fatalf("Failed to read downloaded file: %v", err)
	}
	log.Printf("File content: %s\n", string(content))

	// 6. Remove the object (optional)
	// err = minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	// if err != nil {
	// 	log.Fatalf("Failed to remove object: %v", err)
	// }
	// log.Printf("Successfully removed object: %s\n", objectName)
}
