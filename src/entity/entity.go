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
	Id         int64     `json:"id"`
	Guid       string    `json:"guid"`
	Version    int8      `json:"version"`
	Text       string    `json:"title"`
	UserId     int64     `json:"user_id"`
	CreateDate time.Time `json:"create_date"`
	Deleted    bool      `json:"deleted"`
	Archive    bool      `json:"archive"`
}

type NoteFile struct {
	Id       int64  `json:"id"`
	NoteId   int64  `json:"note_id"`
	Guid     string `json:"guid"`
	Data     []byte `json:"data"`
	Filename string `json:"filename"`
}
