package entity

type NoteService interface {
	SaveNote(note *Note) (*int64, error)

	GetNote(id int64) (*Note, error)

	GetActualNoteByGuid(noteGuid string) (*Note, error)

	DownNoteVersion(noteGuid string) error

	UpNoteVersion(noteGuid string) error

	/* Получить список актуальных записей юзера */
	GetUserActualNotes(userId int64) ([]Note, error)

	/* Получить основную информацию по последним изменённым заметкам */
	GetLastUserNoteMainInfo(userId int64, maxCount int64) ([]MainNoteInfo, error)

	GetNoteVersionHistory(noteGuid string) ([]NoteVersionInfo, error)

	GetNoteFile(noteFileId int64) (*NoteFile, error)
}
