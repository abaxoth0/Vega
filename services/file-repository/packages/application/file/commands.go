package fileapplication

type CreateFileCommand struct {

}

type UpdateFileCommand struct {

}

type DeleteFilesCommand struct {
	IDs	[]string	`json:"ids"`
}

