package minio

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	FileApplication "vega/packages/application/file"
	StorageConnection "vega/packages/infrastructure/object-storage/connection"
)

func connect(driver *Driver) error {
	return driver.Connect(&StorageConnection.Config{
		URL:      "localhost:9000",
		Login:    "minioadmin",
		Password: "minioadmin",
		Token:    "",
		Secure:   false,
	})
}

func TestObjectStorageDriver(t *testing.T) {
	driver := InitDriver()

	t.Run("Connection", func(t *testing.T) {
		t.Log("Checking default driver status...")
		if status := driver.Status(); status != StorageConnection.Disconnected {
			t.Errorf("Invalid driver status, expected \"disconnected\", but got \"%s\"", status.String())
		}
		t.Log("Checking default driver status: OK")

		t.Log("Testing driver.Connect()...")
		if err := connect(driver); err != nil {
			t.Errorf("Connection failed: %v", err)
		}
		t.Log("Testing Driver.Connect(): OK")

		t.Log("Checking updated driver status...")
		if status := driver.Status(); status != StorageConnection.Connected {
			t.Errorf("Invalid driver status, expected \"connected\", but got \"%s\"", status.String())
		}
		t.Log("Checking updated driver status: OK")

		t.Log("Ping connection...")
		if err := driver.Ping(time.Second * 5); err != nil {
			t.Errorf("Failed to ping storage: %v", err)
		}
		t.Log("Ping connection: OK")

		t.Log("Testing Driver.Disconnect()...")
		if err := driver.Disconnect(); err != nil {
			t.Errorf("Failed to gracefully disconnected from storage: %v", err)
		}
		t.Log("Testing Driver.Disconnect(): OK")
	})
}

// Calls handler over each item of iter. Runs each iteration in new goroutine.
// Blocks until all handlers finish their jobs.
func asyncProcess[T any](iter []T, handler func(index int, value T)) {
	wg := new(sync.WaitGroup)
	for i, v := range iter {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler(i, v)
		}()
	}
	time.Sleep(time.Millisecond * 10)
	wg.Wait()
}

func TestUseCasesImplementation(t *testing.T) {
	driver := InitDriver()
	if err := connect(driver); err != nil {
		t.Fatalf("Connection failed")
	}
	bucketName := "vega--auto-test-" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	t.Log("Test bucket name: " + bucketName)

	t.Run("MakeBucket()", func(t *testing.T) {
		t.Log("Testing Driver.MakeBucket()...")
		err := driver.MakeBucket(&FileApplication.MakeBucketCommand{
			Name: bucketName,
		})
		if err != nil {
			t.Errorf("Failed to create a new bucket: %v", err)
		}
		t.Log("Testing Driver.MakeBucket(): OK")

	})

	type testInputs struct {
		path    string
		invalid bool
		empty   bool
		size    int64
	}

	pathOverflow := new(strings.Builder)
	segmentOverflow := new(strings.Builder)
	pathOverflow.WriteByte('/')
	segmentOverflow.WriteByte('/')

	for i := range 1100 {
		pathOverflow.WriteByte('a')

		segmentOverflow.WriteByte('a')
		if i%300 == 0 {
			segmentOverflow.WriteByte('/')
		}
	}

	commonInvalidInputs := []testInputs{
		{path: "", invalid: true},
		{path: "/", invalid: true},
		{path: "//", invalid: true},
		{path: "///", invalid: true},
		{path: ".", invalid: true},
		{path: ".", invalid: true},
		{path: "./", invalid: true},
		{path: pathOverflow.String(), invalid: true},
		{path: segmentOverflow.String(), invalid: true},
	}

	var err error

	t.Run("Mkdir()", func(t *testing.T) {
		dirInputs := append(commonInvalidInputs, []testInputs{
			{path: "/direcotry/"},
			{path: "/direcotry", invalid: true},
		}...)

		for _, input := range dirInputs {
			err = driver.Mkdir(&FileApplication.MkdirCommand{
				Bucket: bucketName,
				Path:   input.path,
			})
			if (err != nil && input.invalid) || (err == nil && !input.invalid) {
				continue
			}
			t.Errorf("Invalid Driver.Mkdir() result. Should fail - %t. Error: %v", input.invalid, err)
		}
	})

	filesPaths := []string{}

	fileContent := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	t.Run("UploadFile()", func(t *testing.T) {
		fileInputs := append(commonInvalidInputs, []testInputs{
			{path: "/test-file.txt"},
			{path: "/f"},
			{path: "/dir/file.txt"},
			{path: "/dir/some.dir/f"},
			{path: "/some-file", empty: true},
			{path: "/some/dir/file.txt", empty: true},
			{path: "/and/another/one", size: 100},
		}...)

		asyncProcess(fileInputs, func(i int, input testInputs) {
			t.Logf("Uploading \"%s\"...", input.path)
			cmd := FileApplication.UploadFileCommand{
				Bucket: bucketName,
				Path:   input.path,
			}
			if !input.empty {
				cmd.Content = strings.NewReader(fileContent)
			} else {
				input.invalid = true
			}
			if input.size == 0 {
				cmd.ContentSize = int64(len(fileContent))
			} else {
				cmd.ContentSize = input.size
				if cmd.ContentSize != int64(len(fileContent)) {
					input.invalid = true
				}
			}
			if !input.invalid {
				filesPaths = append(filesPaths, input.path)
			}
			e := driver.UploadFile(&cmd)
			// Allow invalid inputs to be used, but ignore the result.
			// Just to see will it cause panic or some unexpected behaviour or not.
			if e != nil && !input.invalid {
				t.Errorf("Failed to upload file \"%s\": %v", input.path, e)
			}
		})
	})

	t.Run("GetFileByPath()", func(t *testing.T) {
		asyncProcess(filesPaths, func(_ int, path string) {
			cmd := FileApplication.GetFileByPathQuery{
				Bucket: bucketName,
				Path:   path,
			}
			file, err := driver.GetFileByPath(&cmd)
			if err != nil {
				t.Errorf("Failed to get file \"%s\": %v", path, err)
			}
			if file.Size > 0 {
				content, err := io.ReadAll(file.Content)
				if err != nil {
					t.Errorf("Failed to read content of \"%s\": %v", path, err)
				}
				if string(content) != fileContent {
					t.Errorf("File content doesn't match")
				}
			}
		})
	})

	newFileContent := []byte("some new file content")

	t.Run("UpdateFileContent()", func(t *testing.T) {
		asyncProcess(filesPaths, func(_ int, path string) {
			err = driver.UpdateFileContent(&FileApplication.UpdateFileContentCommand{
				Path:       path,
				Bucket:     bucketName,
				NewContent: bytes.NewReader(newFileContent),
				Size: int64(len(newFileContent)),
			})
			if err != nil {
				t.Fatalf("Failed to update file content \"%s\": %v", path, err)
			}
			file, err := driver.GetFileByPath(&FileApplication.GetFileByPathQuery{
				Bucket: bucketName,
				Path:   path,
			})
			if err != nil {
				t.Fatalf("Failed to get file \"%s\": %v", path, err)
			}
			if file.Size == 0 {
				t.Fatalf("Missing file contnent. File \"%s\"", path)
			}
			content, err := io.ReadAll(file.Content)
			if err != nil {
				t.Fatalf("Failed to read content of \"%s\": %v", path, err)
			}
			if string(content) != string(newFileContent) {
				t.Errorf("New file content doesn't match")
			}
		})
	})

	t.Run("DeleteFiles()", func(t *testing.T) {
		asyncProcess(filesPaths, func(_ int, path string) {
			err = driver.DeleteFiles(&FileApplication.DeleteFilesCommand{
				Bucket: bucketName,
				Paths:  []string{path},
			})
			if err != nil {
				t.Errorf("Failed to delete file \"%s\": %v", path, err)
			}
		})
	})

	t.Run("DeleteBucket()", func(t *testing.T) {
		err = driver.Mkdir(&FileApplication.MkdirCommand{
			Bucket: bucketName,
			Path:   "/prevent-bucket-deletion/"},
		)
		if err != nil {
			t.Errorf("Error: faield to create directory for preventing bucket deletion")
		}
		err = driver.DeleteBucket(&FileApplication.DeleteBucketCommand{
			Name: bucketName,
		})
		if err == nil {
			t.Errorf("Error: non-empty bucket was deleted")
		}
		err = driver.DeleteBucket(&FileApplication.DeleteBucketCommand{
			Name:  bucketName,
			Force: true,
		})
		if err != nil {
			t.Errorf("Failed to delete bucket: %v", err)
		}
	})

	if err := driver.Disconnect(); err != nil {
		t.Logf("Failed to disconnect: %v", err)
	}
}
