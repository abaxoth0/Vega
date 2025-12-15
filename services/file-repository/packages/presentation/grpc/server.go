package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"time"
	objectstorage "vega_file_repository/packages/infrastructure/object-storage"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
)

var ErrServerNotStarted = errors.New("Server is not started, hence can't be stopped.")

type Server struct {
	listening bool
	server    *grpc.Server
	storage   objectstorage.ObjectStorageDriver

	file_repository.UnimplementedFileRepositoryServiceServer
}

func NewServer(storage objectstorage.ObjectStorageDriver) (*Server, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}

	return &Server{storage: storage}, nil
}

func (s *Server) Start(port uint16) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()
	file_repository.RegisterFileRepositoryServiceServer(s.server, s)

	s.listening = true

	if err := s.server.Serve(listener); err != nil {
		s.listening = false
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	if !s.listening {
		return ErrServerNotStarted
	}
	s.server.Stop()
	return nil
}

func (s *Server) HealthCheck(
	ctx context.Context,
	req *file_repository.HealthCheckRequest,
) (*file_repository.HealthCheckResponse, error) {
	log.Printf("Health check called for service: %s", req.GetService())
	return &file_repository.HealthCheckResponse{
		Status:    "SERVING",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
