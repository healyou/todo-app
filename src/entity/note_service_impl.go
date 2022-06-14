package entity

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"todo/src/db"
	"todo/src/filestorage"

	"github.com/pkg/errors"
)

type NoteServiceImpl struct {
	JdbcTemplate db.JdbcTemplate
	MinioService filestorage.MinioService
}

func (service NoteServiceImpl) SaveNote(note *Note) (*int64, error) {
	var id *int64
	var err error
	if note.Id == nil {
		id, err = service.createNote(note)
	} else {
		id, err = service.updateNote(note)
	}

	if err != nil {
		return nil, errors.Wrap(err, "ошибка сохранения note")
	} else {
		return id, nil
	}
}

func (service NoteServiceImpl) createNote(note *Note) (*int64, error) {
	var createdNoteId *int64

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			insertNoteSql := "INSERT INTO note (note_guid, title, text, user_id)" +
				" VALUES (?,?,?,?)"
			result, err := DB.ExecContext(ctx, insertNoteSql, note.NoteGuid, note.Title, note.Text, note.UserId)
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
		return nil, errors.Wrap(err, "ошибка создания записи note")
	} else {
		return createdNoteId, nil
	}
}

func (service NoteServiceImpl) createNoteFile(DB *sql.Tx, ctx context.Context, file *NoteFile, noteId int64) error {
	if file.Data == nil {
		return errors.New("нет данных файла для сохранения")
	}
	if file.Filename == nil {
		return errors.New("не указано имя файла")
	}

	file.Guid = new(string)
	saveFileGuid, err := service.MinioService.SaveFile(file.Data, *file.Filename)
	if err != nil {
		return errors.Wrap(err, "ошибка сохранения файла в minio")
	}
	file.Guid = saveFileGuid

	insertFileSql :=
		"INSERT INTO note_file (note_id, file_guid, filename) " +
			"VALUES (?,?,?)"
	result, err := DB.ExecContext(ctx, insertFileSql, noteId, file.Guid, file.Filename)
	if err != nil {
		return errors.Wrap(err, "ошибка создания записи notefile")
	}
	_, err = result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "ошибка создания записи note")
	}
	return nil
}

func (service NoteServiceImpl) updateNote(note *Note) (*int64, error) {
	var createdNoteId *int64

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			prevNote, err := service.GetActualNoteByGuid(*note.NoteGuid)
			if err != nil {
				return err
			}

			if *prevNote.Id != *note.Id {
				return errors.New("нельзя изменить неактуальную версию note")
			}

			setPrevVersionNotActual(DB, ctx, *prevNote.NoteGuid)
			newNoteId, err := saveNewNoteVersion(DB, ctx, note)
			if err != nil {
				return err
			}
			if newNoteId != nil {
				/* Перемещаем файлы в новую версию */
				err = updateNoteFilesNoteId(DB, ctx, *prevNote.Id, *newNoteId)
				if err != nil {
					return err
				}

				/* Идентификаторы удалённых файлов*/
				var removedFileIds []int64
				removedFileIds, err = getRemovedFileIds(DB, ctx, newNoteId, note.NoteFiles)
				if err != nil {
					return err
				}
				err = service.removeNoteFiles(DB, ctx, newNoteId, removedFileIds)
				if err != nil {
					return err
				}

				/* Новые файлы */
				var newFiles []NoteFile = getNewFiles(note.NoteFiles)
				for i := 0; i < len(newFiles); i++ {
					file := newFiles[i]
					err := service.createNoteFile(DB, ctx, &file, *newNoteId)
					if err != nil {
						return err
					}
				}

				/* Файлы с обновлённым контентом */
				var updatedFiles []NoteFile = getUpdatedFiles(note.NoteFiles)
				for i := 0; i < len(updatedFiles); i++ {
					file := updatedFiles[i]
					err := service.updateNoteFile(DB, ctx, &file, *newNoteId)
					if err != nil {
						return err
					}
				}
			} else {
				return errors.New("не удалось создать note")
			}

			createdNoteId = newNoteId
			return nil
		})

	if err != nil {
		return nil, errors.Wrap(err, "ошибка обновления записи note")
	} else {
		return createdNoteId, nil
	}
}

func (service NoteServiceImpl) updateNoteFile(DB *sql.Tx, ctx context.Context, noteFile *NoteFile, newNoteId int64) error {
	fileGuid, err := service.MinioService.SaveFile(noteFile.Data, *noteFile.Filename)
	if err != nil {
		return errors.Wrap(err, "ошибка сохранения файла в minio")
	}

	updateFileSql := "update note_file set filename = ?, file_guid = ?, note_id = ? where id = ?"
	_, err = DB.ExecContext(ctx, updateFileSql, *noteFile.Filename, fileGuid, newNoteId, noteFile.Id)
	return errors.Wrap(err, "ошибка обновления записи note_file")
}

func getUpdatedFiles(noteFile []NoteFile) []NoteFile {
	var updatedFiles []NoteFile

	for i := 0; i < len(noteFile); i++ {
		file := noteFile[i]
		if file.Id != nil && file.Data != nil {
			updatedFiles = append(updatedFiles, file)
		}
	}

	return updatedFiles
}

func getNewFiles(noteFiles []NoteFile) []NoteFile {
	var newFiles []NoteFile

	for i := 0; i < len(noteFiles); i++ {
		file := noteFiles[i]
		if file.Id == nil {
			newFiles = append(newFiles, file)
		}
	}

	return newFiles
}

func (service NoteServiceImpl) removeNoteFiles(DB *sql.Tx, ctx context.Context, newNoteId *int64, removedFileIds []int64) error {
	if len(removedFileIds) == 1 {
		return nil
	}

	for i := 0; i < len(removedFileIds); i++ {
		fileGuidSql := "select file_guid from note_file where note_id = ? and id = ?"
		row := DB.QueryRowContext(ctx, fileGuidSql, *newNoteId, removedFileIds[i])
		if row.Err() != nil {
			return errors.Wrap(row.Err(), "ошибка получения гуида note")
		}
		var fileGuid string
		row.Scan(&fileGuid)

		removeFileSql := "delete from note_file where note_id = ? and id = ?"
		_, err := DB.ExecContext(ctx, removeFileSql, *newNoteId, removedFileIds[i])
		if err != nil {
			return errors.Wrap(err, "ошибка удаления неактуальных note_file")
		}

		err = service.MinioService.RemoveFile(fileGuid)
		if err != nil {
			return errors.Wrap(err, "ошибка удаления файла из minio")
		}
	}
	return nil
}

func getRemovedFileIds(DB *sql.Tx, ctx context.Context, newNoteId *int64, noteFiles []NoteFile) ([]int64, error) {
	noteFilesSql := "select id from note_file where note_id = ?"
	result, err := DB.QueryContext(ctx, noteFilesSql, *newNoteId)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения id note_file")
	}
	defer result.Close()

	var currentDbFileIds []int64
	for result.Next() {
		var id int64
		if err := result.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "ошибка получения id note_file")
		}
		currentDbFileIds = append(currentDbFileIds, id)
	}

	var removedFileIds []int64
	var noteFileIds []int64
	for i := 0; i < len(noteFiles); i++ {
		file := noteFiles[i]
		if file.Id != nil {
			noteFileIds = append(noteFileIds, *file.Id)
		}
	}

	for i := 0; i < len(currentDbFileIds); i++ {
		dbFileId := currentDbFileIds[i]
		if !intInSlice(dbFileId, noteFileIds) {
			removedFileIds = append(removedFileIds, dbFileId)
		}
	}

	return removedFileIds, nil
}

func intInSlice(a int64, list []int64) bool {
	for b := range list {
		if list[b] == a {
			return true
		}
	}
	return false
}

func updateNoteFilesNoteId(DB *sql.Tx, ctx context.Context, prevNoteId int64, newNoteId int64) error {
	sql := "update note_file set note_id = ? where note_id = ?"
	_, err := DB.ExecContext(ctx, sql, newNoteId, prevNoteId)
	return errors.Wrap(err, "ошибка обновления записи note_file")
}

func saveNewNoteVersion(DB *sql.Tx, ctx context.Context, note *Note) (*int64, error) {
	getMaxNoteVersionSql := "select max(version) as noteMaxVersionNumber from note where note_guid = ?"
	var noteMaxVersionNumber *int8
	row := DB.QueryRowContext(ctx, getMaxNoteVersionSql, *note.NoteGuid)
	if row.Err() != nil {
		return nil, errors.Wrap(row.Err(), "ошибка получения максимальной версии note_file")
	}
	err := row.Scan(&noteMaxVersionNumber)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения версии note")
	}

	insertSql := "INSERT INTO note (prev_note_version_id, note_guid, version, title, text, actual, user_id) VALUES (?, ?, ?, ?, ?, ?, ?)"
	insertResult, err := DB.ExecContext(ctx, insertSql,
		*note.Id, *note.NoteGuid, *noteMaxVersionNumber+1, *note.Title, *note.Text, 1, note.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка сохранения записи note")
	}
	newNoteId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "ошибка сохранения записи note")
	}
	return &newNoteId, nil
}

func setPrevVersionNotActual(DB *sql.Tx, ctx context.Context, noteGuid string) error {
	updateSql := "UPDATE note set actual = 0 where note_guid = ? and actual = 1"
	_, err := DB.ExecContext(ctx, updateSql, noteGuid)
	if err != nil {
		return errors.Wrap(err, "ошибка обновления note")
	}
	return nil
}

func (service NoteServiceImpl) GetActualNoteByGuid(noteGuid string) (*Note, error) {
	var note *Note

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			noteSql := "select * from note where note_guid = ? and actual = 1"
			noteResult := DB.QueryRowContext(ctx, noteSql, noteGuid)
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
		return nil, errors.Wrap(err, "ошибка получения записи note")
	} else {
		return note, nil
	}
}

func (service NoteServiceImpl) GetNote(id int64) (*Note, error) {
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
		return nil, errors.Wrap(err, "ошибка получения записи note")
	} else {
		return note, nil
	}
}

func (service NoteServiceImpl) GetNoteFile(noteFileId int64) (*NoteFile, error) {
	var noteFile *NoteFile

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			noteSql := "select * from note_file where id = ?"
			noteFileResult := DB.QueryRowContext(ctx, noteSql, noteFileId)
			err := noteFileResult.Err()
			if err != nil {
				return err
			}

			noteFile, err = MapNoteFile(noteFileResult)
			if err != nil {
				return err
			}

			data, err := service.MinioService.GetFile(*noteFile.Guid)
			if err != nil {
				return err
			}
			noteFile.Data = data

			return nil
		})

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения записи note_file")
	} else {
		return noteFile, nil
	}
}

func (service NoteServiceImpl) getNoteFilesByNoteId(DB *sql.Tx, ctx context.Context, noteId int64) ([]NoteFile, error) {
	noteFilesSql :=
		"select * from note_file where note_id in (" +
			"    select id from note where id = ?" +
			")"
	noteFilesResult, err := DB.QueryContext(ctx, noteFilesSql, noteId)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка чтения note_file")
	}
	defer func(noteFilesResult *sql.Rows) {
		closeError := noteFilesResult.Close()
		if closeError != nil {
			closeError = errors.Wrap(closeError, "ошибка закрытия resultset")
			log.Println(fmt.Printf("%+v", closeError))
			if err == nil {
				err = closeError
			}
		}
	}(noteFilesResult)
	noteFiles, err := MapNoteFiles(noteFilesResult)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка маппинга note_file")
	}

	return noteFiles, nil
}

func (service NoteServiceImpl) DownNoteVersion(noteGuid string) error {
	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			var prevNoteVersionId, err = getPrevNoteVersionIdIfExists(ctx, DB, noteGuid)
			if err != nil {
				return err
			}

			if prevNoteVersionId != nil {
				setNoteNotActualSql := "update note set actual = 0 where note_guid = ? and actual = 1"
				_, err := DB.ExecContext(ctx, setNoteNotActualSql, noteGuid)
				if err != nil {
					return err
				}

				setPrevNoteActualSql := "update note set actual = 1 where note_guid = ? and id = ?"
				_, err = DB.ExecContext(ctx, setPrevNoteActualSql, noteGuid, *prevNoteVersionId)
				if err != nil {
					return err
				}

				updateNoteFileNoteIdSql := "update note_file set note_id = (select id from note where note_guid = ? and id = ?) where note_id = (select id from note where note_guid = ? and prev_note_version_id = ?)"
				_, err = DB.ExecContext(ctx, updateNoteFileNoteIdSql, noteGuid, *prevNoteVersionId, noteGuid, *prevNoteVersionId)
				if err != nil {
					return err
				}
			} else {
				return errors.New("нельзя уменьшить версию note")
			}

			return nil
		})

	if err != nil {
		return errors.Wrap(err, "ошибка уменьшения версии note")
	} else {
		return nil
	}
}

func (service NoteServiceImpl) UpNoteVersion(noteGuid string) error {
	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			currentActualNoteId, err := getCurrentActualNoteIdByGuid(ctx, DB, noteGuid)
			if err != nil {
				return err
			}
			upNoteVersionId, err := getUpNoteVersionIdIfExists(ctx, DB, *currentActualNoteId)
			if err != nil {
				return err
			}

			if upNoteVersionId != nil {
				setNoteNotActualSql := "update note set actual = 0 where note_guid = ? and actual = 1"
				_, err := DB.ExecContext(ctx, setNoteNotActualSql, noteGuid)
				if err != nil {
					return err
				}

				setPrevNoteActualSql := "update note set actual = 1 where id = ?"
				_, err = DB.ExecContext(ctx, setPrevNoteActualSql, *upNoteVersionId)
				if err != nil {
					return err
				}

				updateNoteFileNoteIdSql := "update note_file set note_id = ? where note_id = ?"
				_, err = DB.ExecContext(ctx, updateNoteFileNoteIdSql, *upNoteVersionId, *currentActualNoteId)
				if err != nil {
					return err
				}
			} else {
				return errors.New("нельзя увеличить версию note")
			}

			return nil
		})

	if err != nil {
		return errors.Wrap(err, "ошибка увеличения версии note")
	} else {
		return nil
	}
}

func getPrevNoteVersionIdIfExists(ctx context.Context, DB *sql.Tx, noteGuid string) (*int64, error) {
	prevNoteVersionIdSql := "select prev_note_version_id from note where note_guid = ? and actual = 1"
	row := DB.QueryRowContext(ctx, prevNoteVersionIdSql, noteGuid)
	if row.Err() != nil {
		return nil, errors.Wrap(row.Err(), "ошибка получения версии note")
	}

	var prevNoteVersionId *int64
	row.Scan(&prevNoteVersionId)

	return prevNoteVersionId, nil
}

/* Получить идентификатор актуальной версии note по гуиду */
func getCurrentActualNoteIdByGuid(ctx context.Context, DB *sql.Tx, noteGuid string) (*int64, error) {
	actualNoteIdVersionSql := "select id from note where note_guid = ? and actual = 1"
	row := DB.QueryRowContext(ctx, actualNoteIdVersionSql, noteGuid)
	if row.Err() != nil {
		return nil, errors.Wrap(row.Err(), "ошибка увеличения id note")
	}

	var id int64
	row.Scan(&id)

	return &id, nil
}

/* Получить идентификатор записи для повышения версии, если такая запись есть */
func getUpNoteVersionIdIfExists(ctx context.Context, DB *sql.Tx, noteId int64) (*int64, error) {
	currentNoteVersionSql := "select id from note where prev_note_version_id = ? and actual = 0"
	row := DB.QueryRowContext(ctx, currentNoteVersionSql, noteId)
	if row.Err() != nil {
		return nil, errors.Wrap(row.Err(), "ошибка увеличения версии note")
	}

	var id *int64
	row.Scan(&id)

	return id, nil
}

func (service NoteServiceImpl) GetUserActualNotes(userId int64) ([]Note, error) {
	var notes []Note

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			/* Получаем список записей */
			notesSql := "select * from note where actual = 1 and user_id = ?"
			userNotesResult, err := DB.QueryContext(ctx, notesSql, userId)

			if err != nil {
				return err
			}
			defer func(userNotesResult *sql.Rows) {
				closeError := userNotesResult.Close()
				if closeError != nil {
					err = closeError
				}
			}(userNotesResult)

			notes, err = MapNotes(userNotesResult)
			if err != nil {
				return err
			}

			/* Получаем файлы для записей */
			for index, note := range notes {
				note.NoteFiles, err = service.getNoteFilesByNoteId(DB, ctx, *note.Id)
				if err != nil {
					return err
				}
				notes[index] = note
			}

			return nil
		})

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения списка note пользователя")
	}

	return notes, nil
}

func (service NoteServiceImpl) GetLastUserNoteMainInfo(userId int64, maxCount int64) ([]MainNoteInfo, error) {
	if maxCount < 1 {
		return nil, errors.New("Нельзя выбрать меньше 1 записи")
	}

	var mainNotesInfo []MainNoteInfo

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			/* Получаем список записей */
			mainNotesInfoSql := "select id,note_guid,version,title,user_id,create_date,actual from note where actual = 1 and user_id = ? order by create_date desc LIMIT ?"
			mainNotesInfoResult, err := DB.QueryContext(ctx, mainNotesInfoSql, userId, maxCount)

			if err != nil {
				return err
			}
			defer func(mainNotesInfoResult *sql.Rows) {
				closeError := mainNotesInfoResult.Close()
				if closeError != nil {
					err = closeError
				}
			}(mainNotesInfoResult)

			mainNotesInfo, err = MapMainNotesInfo(mainNotesInfoResult)
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения списка main note info пользователя")
	}

	return mainNotesInfo, nil
}

func (service NoteServiceImpl) GetNoteVersionHistory(noteGuid string) ([]NoteVersionInfo, error) {
	var noteVersionHistory []NoteVersionInfo

	err := service.JdbcTemplate.ExecuteInTransaction(
		func(ctx context.Context, DB *sql.Tx) error {
			/* Получаем список записей */
			noteHistorySql := "select id as note_id, prev_note_version_id, version, create_date, actual " +
				"from note where note_guid = ? order by create_date desc, id desc"
			noteHistoryResult, err := DB.QueryContext(ctx, noteHistorySql, noteGuid)

			if err != nil {
				return err
			}
			defer func(noteHistoryResult *sql.Rows) {
				closeError := noteHistoryResult.Close()
				if closeError != nil {
					err = closeError
				}
			}(noteHistoryResult)

			noteVersionHistory, err = MapNoteVersionHistoryItems(noteHistoryResult)
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения списка истории изменения note")
	}

	return noteVersionHistory, nil
}
