package grpc

import (
	"context"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Mkdir(context.Context, *file_repository.MkdirRequest) (*emptypb.Empty, error) {
	panic("Mkdir() is not implemented")
}

func (s *Server) UploadFile(grpc.ClientStreamingServer[file_repository.UploadFileRequest, file_repository.UploadFileResponse]) error {
	panic("UploadFile() is not implemented")
}

func (s *Server) UpdateFileContent(grpc.ClientStreamingServer[file_repository.UpdateFileContentRequest, emptypb.Empty]) error {
	panic("UpdateFileContent() is not implemented")
}

func (s *Server) UpdateFileMetadata(context.Context, *file_repository.UpdateFileMetadataRequest) (*emptypb.Empty, error) {
	panic("UpdateFileMetadata() is not implemented")
}

func (s *Server) DeleteFiles(context.Context, *file_repository.DeleteFilesRequest) (*emptypb.Empty, error) {
	panic("DeleteFiles() is not implemented")
}
