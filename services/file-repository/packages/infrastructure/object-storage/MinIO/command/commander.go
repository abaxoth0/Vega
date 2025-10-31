package miniocommand

import FileApplication "vega/packages/application/file"

var Handler FileApplication.CommandHandler = new(defaultCommandHandler)

type defaultCommandHandler struct {

}

func (h *defaultCommandHandler) CreateFile(cmd *FileApplication.CreateFileCommand) error {
	return nil
}

func (h *defaultCommandHandler) UpdateFile(cmd *FileApplication.UpdateFileCommand) error {
	return nil
}

func (h *defaultCommandHandler) DeleteFiles(cmd *FileApplication.DeleteFilesCommand) error {
	return nil
}

