package grpc

import (
	"io"
	"log"
	FileApplication "vega_file_repository/packages/application/file"

	file_repository "github.com/abaxoth0/Vega/common/protobuf/generated/go/services/file-repository"
	"google.golang.org/grpc"
)

const downloadChunkSize int64 = 64 * 1024

func (s *Server) GetFileByPath(
	req *file_repository.GetFileByPathRequest,
	stream grpc.ServerStreamingServer[file_repository.FileChunk],
) error {
	fileStream, err := s.storage.GetFileByPath(&FileApplication.GetFileByPathQuery{
		Bucket: req.GetBucket(),
		Path:   req.GetPath(),
	})
	if err != nil {
		return err
	}
	defer fileStream.Cancel()

	chunkSize := downloadChunkSize
	if req.ChunkSize > 0 {
		chunkSize = int64(req.ChunkSize)
	}

	buf := make([]byte, chunkSize)
	var chunkIndex int64
	var totalChunks int64

	if fileStream.Size == 0 {
		totalChunks = 1
	} else {
		totalChunks = (fileStream.Size-1)/chunkSize + 1
	}

	log.Printf(
		"Sending file \"%s\": file size %d bytes; total chunks %d; chunk size %d\n",
		req.GetPath(), fileStream.Size, totalChunks, chunkSize,
	)

	streaming := true
	for streaming {
		n, err := fileStream.Content.Read(buf)
		if err == io.EOF {
			streaming = false
		} else if err != nil {
			return err
		}

		chunk := &file_repository.FileChunk{
			Content:    buf[:n],
			ChunkIndex: chunkIndex,
			TotalSize:  fileStream.Size,
		}
		if err := stream.Send(chunk); err != nil {
			return err
		}
		chunkIndex++
	}

	return nil
}
