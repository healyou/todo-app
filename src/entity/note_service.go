package entity

type NoteService interface {
	SaveNote(note *Note) (*int64, error)

	GetNote(id int64) (*Note, error)

	GetActualNoteByGuid(noteGuid string) (*Note, error)

	DownNoteVersion(noteGuid string) error

	UpNoteVersion(noteGuid string) error

	// TODO добавить метод
	// GetUserNotes(userId int64) ([]Note, error)
}
