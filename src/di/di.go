package di

import (
	"todo/src/entity"
	"todo/src/filestorage"
)

/* Класс для классов одиночек */
type DependencyInjection interface {
	GetNoteService() entity.NoteService
	GetMinioService() filestorage.MinioService
}
