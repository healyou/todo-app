package entity

type NoteService interface {
	Test() error

	SaveNote(note *Note) (*int64, error)

	GetNote(id int64) (*Note, error)

	GetNoteByGuid(noteGuid string) (*Note, error)

	DownNoteVersion(noteGuid string) error

	UpNoteVersion(noteGuid string) error
}
