package main

import (
	"bytes"
	"context"
	"log"
	"time"
	fileapplication "vega/packages/application/file"
	ObjectStorage "vega/packages/infrastructure/object-storage"
	StorageConnection "vega/packages/infrastructure/object-storage/connection"
	"vega/packages/presentation/grpc"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	testStorage()
	// testgRPC()
}

func testgRPC() {
	err := ObjectStorage.Driver.Connect(&StorageConnection.Config{
		URL:      "localhost:9000",
		Login:    "minioadmin",
		Password: "minioadmin",
		Token:    "",
		Secure:   false,
	})
	if err != nil {
		panic(err)
	}
	defer ObjectStorage.Driver.Disconnect()

	server, err := grpc.NewServer(ObjectStorage.Driver)
	if err != nil {
		panic(err)
	}

	println("gRPC server started on port 50001")

	if err := server.Start(50001); err != nil {
		panic(err)
	}
}

func testStorage() {
	err := ObjectStorage.Driver.Connect(&StorageConnection.Config{
		URL:      "localhost:9000",
		Login:    "minioadmin",
		Password: "minioadmin",
		Token:    "",
		Secure:   false,
	})
	if err != nil {
		panic(err)
	}

	if err := ObjectStorage.Driver.Ping(time.Second * 5); err != nil {
		panic("failed to ping object storage")
	}

	_, err = ObjectStorage.Driver.GetFileByPath(&fileapplication.GetFileByPathQuery{
		Path:   "/newfile.txt",
		Bucket: "test-bucket",
	})
	if err != nil {
		panic(err)
	}
	// data, err := io.ReadAll(file.Content)
	// if err != nil {
	// 	panic(err)
	// }
	// println(string(data))

	// e := ObjectStorage.Driver.UploadFile(&fileapplication.UploadFileCommand{
	// 	FileMeta:    nil,
	// 	Content:     file.Content,
	// 	ContentSize: file.Size,
	// 	Path:        "/file.txt",
	// 	Bucket:      "test-bucket",
	// })
	// e := objectstorage.Driver.DeleteFiles(&fileapplication.DeleteFilesCommand{
	// 	Paths: []string{"/test1.txt", "/test2.txt", "/test4.txt"},
	// 	Bucket: "test-bucket",
	// })
	e := ObjectStorage.Driver.UpdateFileContent(&fileapplication.UpdateFileContentCommand{
		Path:       "/my-new-file.txt",
		Bucket:     "test-bucket",
		NewContent: bytes.NewReader([]byte("full replace upd test")),
	})

	// meta, e := objectstorage.Driver.GetFileMetadataByPath(&fileapplication.GetFileByPathQuery{
	// 	Path: "/my-new-file.txt",
	// 	Bucket: "test-bucket",
	// })
	//

	// e := objectstorage.Driver.Mkdir(&fileapplication.MkdirCommand{
	// 	Path: "/my/directory/",
	// 	Bucket: "test-bucket",
	// })
	if e != nil {
		panic(e)
	}

	// fmt.Printf("%+v\n", *meta)

	println("OK")
}

func test() {
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
	bucketName := "test-bucket"
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
	objectName := "war-and-peace.txt"
	filePath := "./war-and-peace.txt"

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
	// downloadPath := "./downloaded-file.txt"
	// err = minioClient.FGetObject(ctx, bucketName, objectName, downloadPath, minio.GetObjectOptions{})
	// if err != nil {
	// 	log.Fatalf("Failed to download file: %v", err)
	// }
	// log.Printf("Successfully downloaded file to: %s\n", downloadPath)

	// 5. Check if file exists and read its content
	// content, err := os.ReadFile(downloadPath)
	// if err != nil {
	// 	log.Fatalf("Failed to read downloaded file: %v", err)
	// }
	// log.Printf("File content: %s\n", string(content))

	// 6. Remove the object (optional)
	// err = minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	// if err != nil {
	// 	log.Fatalf("Failed to remove object: %v", err)
	// }
	// log.Printf("Successfully removed object: %s\n", objectName)
}
