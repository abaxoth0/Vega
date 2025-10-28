package miniocommandhandler

import FileApplication "vega/packages/application/file"

func Init() FileApplication.CommandHandler {
	 return new(defaultCommandHandler)
}

type defaultCommandHandler struct {

}

func (h *defaultCommandHandler) CreateFile(cmd FileApplication.CreateFileCommand) error {
	return nil
}

func (h *defaultCommandHandler) UpdateFile(cmd FileApplication.UpdateFileCommand) error {
	return nil
}

func (h *defaultCommandHandler) DeleteFile(cmd FileApplication.DeleteFileCommand) error {
	return nil
}

