package entity

import "time"

// Album represents data about a record album.
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type Note struct {
	Id                *int64     `json:"id"`
	PrevNoteVersionId *int64     `json:"prev_note_version_id"`
	NoteGuid          *string    `json:"guid"`
	Version           *int8      `json:"version"`
	Title             *string    `json:"title"`
	Text              *string    `json:"text"`
	UserId            *int64     `json:"user_id"`
	CreateDate        *time.Time `json:"create_date"`
	Deleted           *bool      `json:"deleted"`
	Archive           *bool      `json:"archive"`
	Actual            *bool      `json:"actual"`
	NoteFiles         []NoteFile `json:"note_files"`
}

type MainNoteInfo struct {
	Id                *int64     `json:"id"`
	NoteGuid          *string    `json:"guid"`
	Version           *int8      `json:"version"`
	Title             *string    `json:"title"`
	UserId            *int64     `json:"user_id"`
	CreateDate        *time.Time `json:"create_date"`
	Actual            *bool      `json:"actual"`
}

type NoteFile struct {
	Id       *int64  `json:"id"`
	NoteId   *int64  `json:"note_id"`
	Guid     *string `json:"guid"`
	Data     []byte  `json:"data"`
	Filename *string `json:"filename"`
}

type NoteVersionInfo struct {
	NoteId            *int64     `json:"note_id"`
	PrevNoteVersionId *int64     `json:"prev_note_version_id"`
	Version           *int8      `json:"version"`
	CreateDate        *time.Time `json:"create_date"`
	Actual            *bool      `json:"actual"`
}
