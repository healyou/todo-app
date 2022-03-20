package filestorage

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

const minioFilenameUserOptionName = "Filename"
const minioBucketName = "todo-app-bucket"

type MinioServiceImpl struct {
	Client *minio.Client
}

func (service MinioServiceImpl) SaveFile(data []byte, filename string) (*string, error) {
	uuidWithHyphen := uuid.New().String()
	mimeType := mimetype.Detect(data)

	userMetaData := map[string]string{
		minioFilenameUserOptionName: filename,
	}
	options := minio.PutObjectOptions{
		ContentType:  mimeType.String(),
		UserMetadata: userMetaData,
	}
	ctx := context.Background()
	defer ctx.Done()
	info, err := service.Client.PutObject(
		ctx, minioBucketName, uuidWithHyphen, bytes.NewReader(data), int64(len(data)), options)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка сохранения объекта в minio")
	}

	return &info.Key, nil
}

func (service MinioServiceImpl) GetFile(fileUuid string) ([]byte, error) {
	ctx := context.Background()
	defer ctx.Done()

	object, err := service.Client.GetObject(ctx, minioBucketName, fileUuid, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения объекта из minio")
	}
	defer func(object *minio.Object) {
		closeErr := object.Close()
		if closeErr != nil {
			closeErr = errors.Wrap(closeErr, "ошибка закрытия объекта minio после работы с ним")
			log.Println(fmt.Printf("%+v", closeErr))
			if (err == nil) {
				err = closeErr
			}
		}
	}(object)

	objectInfo, err := object.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения информации об объекте minio")
	}

	dataBuffer := new(bytes.Buffer)
	copySize, err := dataBuffer.ReadFrom(object)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка копирования файла")
	}
	if copySize != objectInfo.Size {
		return nil, errors.New("число скопированнх данных не совпадает с нужным числом")
	}

	return dataBuffer.Bytes(), err
}

func (service MinioServiceImpl) RemoveFile(fileUuid string) error {
	ctx := context.Background()
	defer ctx.Done()
	return service.Client.RemoveObject(ctx, minioBucketName, fileUuid, minio.RemoveObjectOptions{})
}
