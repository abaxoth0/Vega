package grpc

import (
	"context"
	"io"
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
		t.Logf("gRPC server started")
		if err := server.Start(testPort); err != nil {
			t.Fatalf("Failed to start gRPC server: %v", err)
			return
		}
		t.Logf("gRPC server stopped")
	}()

	time.Sleep(time.Millisecond * 20)
	defer func() {
		t.Logf("Stopping gRPC server...")
		if err := server.Stop(); err != nil {
			t.Fatalf("Failed to start gRPC server: %v", err)
			return
		}
		t.Logf("Stopping gRPC server: DONE")
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

	t.Run("GetFileByPath()", func(t *testing.T) {
		withClient(t, func(client file_repository.FileRepositoryServiceClient) {
			ctx, cancel := newRPCContext()
			defer cancel()

			fileStream, err := client.GetFileByPath(ctx, &file_repository.GetFileByPathRequest{
				Bucket: "test-bucket",
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
}
