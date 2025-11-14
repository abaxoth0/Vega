package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"time"
	objectstorage "vega/packages/infrastructure/object-storage"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
)

type server struct {
	storage objectstorage.ObjectStorageDriver
	file_repository.UnimplementedFileRepositoryServiceServer
}

func NewServer(storage objectstorage.ObjectStorageDriver) (*server, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}

	return &server{storage: storage}, nil
}

func (s *server) Start(port uint16) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	file_repository.RegisterFileRepositoryServiceServer(grpcServer, s)

	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *server) HealthCheck(
	ctx context.Context,
	req *file_repository.HealthCheckRequest,
) (*file_repository.HealthCheckResponse, error){
	log.Printf("Health check called for service: %s", req.GetService())
	return &file_repository.HealthCheckResponse{
		Status: "SERVING",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
