package main

import (
	"fmt"
	"time"
	"vega_file_discovery/cmd/app"
	"vega_file_discovery/common/config"
	fileapplication "vega_file_discovery/packages/application/file"
	"vega_file_discovery/packages/entity"
	DB "vega_file_discovery/packages/infrastrcuture/database"

	"github.com/abaxoth0/Vega/libs/go/packages/logger"
)

var log = logger.NewSource("MAIN", logger.Default)

func main() {
	app.StartInit()
		app.InitDefault()

		logger.Debug.Store(config.Debug.Enabled)
		logger.Trace.Store(config.App.TraceLogsEnabled)
	app.EndInit()

	go func() {
		if err := logger.Default.Start(config.Debug.Enabled); err != nil {
			panic(err.Error())
		}
	}()
	defer func() {
		if err := logger.Default.Stop(); err != nil {
			log.Error("Failed to stop logger", err.Error(), nil)
		}
	}()

	// Reserve some time for logger to start up
	time.Sleep(time.Millisecond * 50)

	if err := DB.Database.Connect(); err != nil {
		panic(err)
	}

	fileContent := "some text idk..."

	id, err := DB.Database.CreateFileMetadata(&fileapplication.CreateFileMetadataCmd{
		Metadata: &entity.FileMetadata{
			UpdatableFileMetadata: entity.UpdatableFileMetadata{
				Bucket: "ced077c7-08de-4237-9358-9778780e0592",
				Path:   "/example.txt",
				Encoding: "UTF-8",
				Owner: "cf2dbfb0-deeb-4c3c-9679-946dd31e9dd7",
				Permissions: entity.NewFilePermissions(
					entity.ReadFilePermission|entity.UpdateFilePermission|entity.DeleteFilePermission,
					entity.ReadFilePermission,
					0,
				),
				Categories: []string{"test", "example", "dev"},
				Tags: []string{"testing"},
				Status: entity.ActiveFileStatus,
			},
			GeneratedFileMetadata: entity.GeneratedFileMetadata{
				OriginalName: "example.txt",
				MIMEType: "text/plain",
				Size: int64(len(fileContent)),
				// Checksum: string(entity.ChecksumHasher().Sum([]byte(fileContent))),
				Checksum: entity.HashAll([]byte(fileContent)),
				UploadedBy: "cf2dbfb0-deeb-4c3c-9679-946dd31e9dd7",
			},
		},
	})

	if err != nil {
		fmt.Printf("[ ERROR ] Failed to create metadata: %v\n", err.Error())
	}

	if err := DB.Database.Disconnect(); err != nil {
		panic(err)
	}

	fmt.Printf("Done (file id: %s)\n", id)
	x := ""
	fmt.Scan(&x)
}
