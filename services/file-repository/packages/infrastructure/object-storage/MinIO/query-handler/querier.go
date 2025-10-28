package minioqueryhandler

import (
	FileApplication "vega/packages/application/file"
	"vega/packages/domain/entity"
)

func Init() FileApplication.QueryHandler {
	return new(defaultQueryHandler)
}

type defaultQueryHandler struct {

}

func (m *defaultQueryHandler) GetFileByID(query FileApplication.GetFileByIdQuery) (*entity.File, error) {
	return nil, nil
}

