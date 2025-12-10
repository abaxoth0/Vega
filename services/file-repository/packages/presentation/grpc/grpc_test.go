package grpc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
	objectstorage "vega/packages/infrastructure/object-storage"
	storageconnection "vega/packages/infrastructure/object-storage/connection"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const testPort uint16 = 50001

type closeFunc = func() error

func new_grpc_clinet() (file_repository.FileRepositoryServiceClient, closeFunc) {
	conn, err := grpc.NewClient(
		"localhost:"+strconv.Itoa(int(testPort)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := file_repository.NewFileRepositoryServiceClient(conn)
	return client, conn.Close
}

func newRPCContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

// Creates gRPC server and clinet, after that pass this client into the handler function.
// Handles connection/disconnection automatically.
func withClient(t *testing.T, handler func(client file_repository.FileRepositoryServiceClient)) {
	server, err := NewServer(objectstorage.Driver)
	if err != nil {
		t.Fatalf("Failed to create gRPC server: %v", err)
		return
	}

	go func() {
		if err := server.Start(testPort); err != nil {
			t.Fatalf("Failed to start gRPC server: %v", err)
			return
		}
	}()

	time.Sleep(time.Millisecond * 20)
	defer func() {
		if err := server.Stop(); err != nil {
			t.Fatalf("Failed to start gRPC server: %v", err)
			return
		}
	}()

	client, closeClient := new_grpc_clinet()
	defer closeClient()

	handler(client)
}

func TestServer(t *testing.T) {
	withClient(t, func(client file_repository.FileRepositoryServiceClient) {
		ctx, cancel := newRPCContext()
		defer cancel()

		healthResp, err := client.HealthCheck(ctx, &file_repository.HealthCheckRequest{
			Service: "file-repository",
		})
		if err != nil {
			t.Fatalf("HealthCheck RPC failed: %v", err)
		}

		t.Logf("Health check status: %s", healthResp.GetStatus())
		t.Logf("Health check timestamp: %s", healthResp.GetTimestamp())
	})
}

func TestRPC(t *testing.T) {
	err := objectstorage.Driver.Connect(&storageconnection.Config{
		URL:      "localhost:9000",
		Login:    "minioadmin",
		Password: "minioadmin",
		Token:    "",
		Secure:   false,
	})
	if err != nil {
		t.Fatalf("Failed to connect to object storage: %v", err)
	}
	defer func() {
		if err := objectstorage.Driver.Disconnect(); err != nil {
			t.Fatalf("Failed to disconnect from object storage")
		}
	}()

	const testBucket string = "test-bucket"

	t.Run("GetFileByPath()", func(t *testing.T) {
		withClient(t, func(client file_repository.FileRepositoryServiceClient) {
			ctx, cancel := newRPCContext()
			defer cancel()

			fileStream, err := client.GetFileByPath(ctx, &file_repository.GetFileByPathRequest{
				Bucket: testBucket,
				Path:   "/file.txt",
			})
			if err != nil {
				t.Fatalf("GetFileByPath() RPC failed: %v", err)
			}
			streaming := true
			for streaming {
				chunk, err := fileStream.Recv()
				if err == io.EOF {
					streaming = false
				} else if err != nil {
					t.Fatalf("File stream failed: %v", err)
					break
				}
				if streaming && chunk == nil {
					t.Fatalf("File stream failed: chunk is nil")
					break
				}
				t.Logf("Chunk content: %v\n", string(chunk.GetContent()))
			}
		})
	})

	t.Run("Mkdir()", func(t *testing.T) {
		withClient(t, func(client file_repository.FileRepositoryServiceClient) {
			context, cancel := newRPCContext()
			defer cancel()

			_, err := client.Mkdir(context, &file_repository.MkdirRequest{
				Bucket: testBucket,
				Path: fmt.Sprintf("/test/mkdir-%s/", time.Now().Format(time.RFC3339)),
			})
			if err != nil {
				t.Fatalf("Mkdir() RPC failed: %v", err)
			}
		})
	})

	t.Run("UploadFile()", func(t *testing.T) {
		withClient(t, func(client file_repository.FileRepositoryServiceClient) {
			fileContent := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

			if err := testFileStream(client.UploadFile, testBucket, fileContent); err != nil {
				t.Fatalf("UpdateFileContent() RPC failed: %v", err)
			}
		})
	})

	t.Run("UpdateFileContent()", func(t *testing.T) {
		withClient(t, func(client file_repository.FileRepositoryServiceClient) {
			err := testFileStream(client.UpdateFileContent, testBucket, []byte("new file content"))

			if err != nil {
				t.Fatalf("UpdateFileContent() RPC failed: %v", err)
			}
		})
	})
}

type fileStreamFunc = func (ctx context.Context, opts ...grpc.CallOption) (
	grpc.BidiStreamingClient[file_repository.FileContentRequest, file_repository.StatusResponse],
	error,
)

func testFileStream(streamFunc fileStreamFunc, bucket string, content []byte) error {
	context, cancel := newRPCContext()
	defer cancel()

	stream, err := streamFunc(context)
	if err != nil {
		return err
	}

	err = stream.Send(&file_repository.FileContentRequest{
		Data: &file_repository.FileContentRequest_Header{
			Header: &file_repository.FileContentHeader{
				Path: fmt.Sprintf("/upload-file-test-%s", time.Now().Format(time.RFC3339)),
				Bucket: bucket,
				Size: int64(len(content)),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to send header")
	}

	err = stream.Send(&file_repository.FileContentRequest{
		Data: &file_repository.FileContentRequest_Chunk{
			Chunk: content,
		},
	})
	if err != nil {
		return fmt.Errorf("Faield to send chunk")
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	if resp.Status != http.StatusOK {
		return fmt.Errorf("Server reported upload failure")
	}

	return nil
}
