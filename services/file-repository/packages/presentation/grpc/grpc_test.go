package grpc

import (
	"context"
	"testing"
	"time"
	objectstorage "vega/packages/infrastructure/object-storage"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const testPort uint16 = 50001

type closeFunc = func() error

func new_grpc_clinet() (file_repository.FileRepositoryServiceClient, closeFunc) {
	conn, err := grpc.NewClient(
		"localhost:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	client := file_repository.NewFileRepositoryServiceClient(conn)
	return client, conn.Close
}

func newRequestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

func TestServer(t *testing.T) {
	server, err := NewServer(objectstorage.Driver)
	if err != nil {
		t.Fatalf("Failed to create gRPC server: %v", err)
	}

	go func(){
		t.Logf("gRPC server started")
		if err := server.Start(testPort); err != nil {
			t.Fatalf("Failed to start gRPC server: %v", err)
		}
		t.Logf("gRPC server stopped")
	}()

	time.Sleep(time.Millisecond * 20)
	defer func(){
		t.Logf("Stopping gRPC server...")
		if err := server.Stop(); err != nil {
			t.Fatalf("Failed to start gRPC server: %v", err)
		}
		t.Logf("Stopping gRPC server: DONE")
	}()

	client, closeClient := new_grpc_clinet()
	defer closeClient()

	ctx, cancel := newRequestContext()
	defer cancel()

	healthResp, err := client.HealthCheck(ctx, &file_repository.HealthCheckRequest{
		Service: "file-repository",
	})
	if err != nil {
		t.Fatalf("Health check request failed: %v", err)
	}

	t.Logf("Health check status: %s", healthResp.GetStatus())
	t.Logf("Health check timestamp: %s", healthResp.GetTimestamp())
}
