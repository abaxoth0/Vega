package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	fileapplication "vega/packages/application/file"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
)

func (s *Server) Mkdir(
	ctx context.Context,
	req *file_repository.MkdirRequest,
) (*file_repository.StatusResponse, error) {
	err := s.storage.Mkdir(&fileapplication.MkdirCommand{
		Bucket: req.GetBucket(),
		Path: req.GetPath(),
	})
	if err != nil {
		return nil, err
	}
	return &file_repository.StatusResponse{
		Status: http.StatusOK,
	}, nil
}

type fileContent struct {
	Header 	*file_repository.FileContentHeader
	Reader io.Reader
}

func fileContentFromStream(
    stream grpc.BidiStreamingServer[file_repository.FileContentRequest, file_repository.StatusResponse],
) (*fileContent, error) {
    firstMsg, err := stream.Recv()
    if err != nil {
        return nil, fmt.Errorf("failed to receive first message: %v", err)
    }

    header := firstMsg.GetHeader()
    if header == nil {
        return nil, errors.New("first message must be a header")
    }

    pr, pw := io.Pipe()

    go func() {
        defer pw.Close()

        if chunk := firstMsg.GetChunk(); chunk != nil {
            if _, err := pw.Write(chunk); err != nil {
                return
            }
        }

        for {
            msg, err := stream.Recv()
            if err == io.EOF {
                return
            }
            if err != nil {
                return
            }

            if chunk := msg.GetChunk(); chunk != nil {
                if _, err := pw.Write(chunk); err != nil {
                    return
                }
            }
        }
    }()

	return &fileContent{
		Header: header,
		Reader: pr,
	}, nil
}

func (s *Server) UploadFile(
    stream grpc.BidiStreamingServer[file_repository.FileContentRequest, file_repository.StatusResponse],
) error {
	content, err := fileContentFromStream(stream)
	if err != nil {
		return err
	}

    err = s.storage.UploadFile(&fileapplication.UploadFileCommand{
        Bucket:      content.Header.Bucket,
        Path:        content.Header.Path,
        ContentSize: content.Header.Size,
        Content:     content.Reader,
    })
    if err != nil {
        return fmt.Errorf("file upload failed: %v", err)
    }

    return stream.Send(&file_repository.StatusResponse{
        Status: http.StatusOK,
    })
}

func (s *Server) UpdateFileContent(
	stream grpc.BidiStreamingServer[file_repository.FileContentRequest, file_repository.StatusResponse],
) error {
	content, err := fileContentFromStream(stream)
	if err != nil {
		return err
	}

    err = s.storage.UpdateFileContent(&fileapplication.UpdateFileContentCommand{
        Bucket:      content.Header.Bucket,
        Path:        content.Header.Path,
		Size: 		 content.Header.Size,
        NewContent:  content.Reader,
    })
    if err != nil {
        return fmt.Errorf("file update failed: %v", err)
    }

    return stream.Send(&file_repository.StatusResponse{
        Status: http.StatusOK,
    })
}

func (s *Server) DeleteFiles(
	context.Context,
	*file_repository.DeleteFilesRequest,
) (*file_repository.StatusResponse, error) {
	panic("DeleteFiles() is not implemented")
}
