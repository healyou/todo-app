package entity

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
	"todo/src/db"
	"todo/src/filestorage"
)

type NoteServiceImpl struct {
	JdbcTemplate db.JdbcTemplate
	MinioService filestorage.MinioService
}

func (service *NoteServiceImpl) Test() error {
	sqlFunc := func(context context.Context, DB *sql.Tx) (*sql.Result, error) {
		sqlCount := "select count(*) from note"
		var cnt int
		err := DB.QueryRowContext(context, sqlCount).Scan(&cnt)
		if err != nil {
			return nil, err
		}

		newUuid := uuid.New().String()
		sqlInsert := "INSERT INTO note (note_guid, text, user_id)\n VALUES ('" + newUuid + "', 'note text1', 1)"
		result, err := DB.ExecContext(context, sqlInsert)
		if err != nil {
			return nil, err
		}

		sqlCount = "select count(*) from note"
		err = DB.QueryRowContext(context, sqlCount).Scan(&cnt)
		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	result, err := service.JdbcTemplate.InTransactionForSqlResult(sqlFunc)
	if err != nil {
		return err
	}
	if result != nil {
		log.Println(result)
	}

	return nil
}

func (service *NoteServiceImpl) SaveNote(note *Note) (*int64, error) {
	var id *int64
	var err error
	if note.Id == nil {
		id, err = service.createNote(note)
	} else {
		id, err = service.updateNote(note)
	}

	if err != nil {
		return nil, err
	} else {
		return id, nil
	}
}

func (service *NoteServiceImpl) createNote(note *Note) (*int64, error) {
	var createdNoteId *int64

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			insertNoteSql := "INSERT INTO note (note_guid, text, user_id)" +
				" VALUES (?,?,?)"
			result, err := DB.ExecContext(ctx, insertNoteSql, note.NoteGuid, note.Text, note.UserId)
			if err != nil {
				return err
			}
			createdNoteId = new(int64)
			*createdNoteId, err = result.LastInsertId()
			if err != nil {
				return err
			}

			if note.NoteFiles != nil && len(note.NoteFiles) > 0 {
				for i := 0; i < len(note.NoteFiles); i++ {
					file := note.NoteFiles[i]
					err := service.createNoteFile(DB, ctx, &file, *createdNoteId)
					if err != nil {
						return err
					}
				}
			}

			return nil
		})

	if err != nil {
		return nil, err
	} else {
		return createdNoteId, nil
	}
}

func (service *NoteServiceImpl) createNoteFile(DB *sql.Tx, ctx context.Context, file *NoteFile, noteId int64) error {
	if file.Data == nil {
		return errors.New("нет данных файла для сохранения")
	}
	if file.Filename == nil {
		return errors.New("не указано имя файла")
	}

	file.Guid = new(string)
	saveFileGuid, err := service.MinioService.SaveFile(file.Data, *file.Filename)
	if err != nil {
		return err
	}
	file.Guid = saveFileGuid

	insertFileSql :=
		"INSERT INTO note_file (note_id, file_guid, filename) " +
			"VALUES (?,?,?)"
	result, err := DB.ExecContext(ctx, insertFileSql, noteId, file.Guid, file.Filename)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (service *NoteServiceImpl) updateNote(note *Note) (*int64, error) {
	var createdNoteId *int64

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			return errors.New("не реализовано обновление")
		})

	if err != nil {
		return nil, err
	} else {
		return createdNoteId, nil
	}
}

func (service *NoteServiceImpl) GetNoteByGuid(noteGuid string) (*Note, error) {
	var note *Note

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			noteSql :=
				"select * from note where note_guid = ? and version = (" +
					"    SELECT MAX(version)" +
					"    FROM note" +
					"    where note_guid = ?" +
					"    GROUP BY note_guid" +
					")"
			noteResult := DB.QueryRowContext(ctx, noteSql, noteGuid, noteGuid)
			if noteResult.Err() != nil {
				return noteResult.Err()
			}
			var err error
			note, err = MapNote(noteResult)
			if err != nil {
				return err
			}

			note.NoteFiles, err = service.getNoteFilesByNoteId(DB, ctx, *note.Id)
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		return nil, err
	} else {
		return note, nil
	}
}

func (service *NoteServiceImpl) GetNote(id int64) (*Note, error) {
	var note *Note

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			noteSql := "select * from note where id = ?"
			noteResult := DB.QueryRowContext(ctx, noteSql, id)
			if noteResult.Err() != nil {
				return noteResult.Err()
			}
			var err error
			note, err = MapNote(noteResult)
			if err != nil {
				return err
			}

			note.NoteFiles, err = service.getNoteFilesByNoteId(DB, ctx, *note.Id)
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		return nil, err
	} else {
		return note, nil
	}
}

func (service *NoteServiceImpl) getNoteFilesByNoteId(DB *sql.Tx, ctx context.Context, noteId int64) ([]NoteFile, error) {
	noteFilesSql :=
		"select * from note_file where note_id in (" +
			"    select id from note where id = ?" +
			")"
	noteFilesResult, err := DB.QueryContext(ctx, noteFilesSql, noteId)
	if err != nil {
		return nil, err
	}
	defer func(noteFilesResult *sql.Rows) {
		closeError := noteFilesResult.Close()
		if closeError != nil {
			err = closeError
		}
	}(noteFilesResult)
	noteFiles, err := MapNoteFiles(noteFilesResult)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(noteFiles); i++ {
		file := &noteFiles[i]
		data, err := service.MinioService.GetFile(*file.Guid)
		if err != nil {
			return nil, err
		}
		file.Data = data
	}

	return noteFiles, nil
}
