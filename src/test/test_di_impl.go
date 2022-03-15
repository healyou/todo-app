package test

import (
	"todo/src/entity"
	"todo/src/filestorage"
)

type TestDependencyInjectionImpl struct {
	NoteServiceValue  entity.NoteService
	MinioServiceValue filestorage.MinioService
}

func (di TestDependencyInjectionImpl) GetNoteService() entity.NoteService {
	return di.NoteServiceValue
}

func (di TestDependencyInjectionImpl) GetMinioService() filestorage.MinioService {
	return di.MinioServiceValue
}
