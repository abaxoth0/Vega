package grpc

import (
	"context"
	fileapplication "vega/packages/application/file"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func empty() *emptypb.Empty {
	return new(emptypb.Empty)
}

func (s *Server) Mkdir(ctx context.Context, req *file_repository.MkdirRequest) (*emptypb.Empty, error) {
	err := s.storage.Mkdir(&fileapplication.MkdirCommand{
		Bucket: req.GetBucket(),
		Path: req.GetPath(),
	})
	if err != nil {
		return nil, err
	}
	return empty(), nil
}

func (s *Server) UploadFile(grpc.ClientStreamingServer[file_repository.UploadFileRequest, file_repository.UploadFileResponse]) error {
	panic("UploadFile() is not implemented")
}

func (s *Server) UpdateFileContent(grpc.ClientStreamingServer[file_repository.UpdateFileContentRequest, emptypb.Empty]) error {
	panic("UpdateFileContent() is not implemented")
}

func (s *Server) DeleteFiles(context.Context, *file_repository.DeleteFilesRequest) (*emptypb.Empty, error) {
	panic("DeleteFiles() is not implemented")
}
