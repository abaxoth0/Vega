package grpc

import (
	"context"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"github.com/abaxoth0/Vega/common/protobuf/generated/go/types"
	"google.golang.org/grpc"
)

func (s *server) GetFileByPath(
	req *file_repository.GetFileByPathRequest,
	stream grpc.ServerStreamingServer[file_repository.FileChunk],
) (error){
	panic("GetFileByPath() is not implemented")
}

func (s *server) GetFileMetadataByPath(context.Context, *file_repository.GetFileMetadataByPathRequest) (*types.FileMetadata, error) {
	panic("GetFileMetadataByPath() is not implemented")
}
