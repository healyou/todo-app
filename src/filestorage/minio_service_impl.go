package filestorage

import (
	"bytes"
	"context"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"log"
)

const MinioFilenameUserOptionName = "Filename"
const MinioBucketName = "todo-app-bucket"

type MinioServiceImpl struct {
	Client *minio.Client
}

func (service MinioServiceImpl) SaveFile(data []byte, filename string) (*string, error) {
	uuidWithHyphen := uuid.New().String()
	mimeType := mimetype.Detect(data)

	userMetaData := map[string]string{
		MinioFilenameUserOptionName: filename,
	}
	options := minio.PutObjectOptions{
		ContentType:  mimeType.String(),
		UserMetadata: userMetaData,
	}
	ctx := context.Background()
	defer ctx.Done()
	info, err := service.Client.PutObject(
		ctx, MinioBucketName, uuidWithHyphen, bytes.NewReader(data), int64(len(data)), options)
	if err != nil {
		return nil, err
	}

	return &info.Key, nil
}

func (service MinioServiceImpl) GetFile(fileUuid string) ([]byte, error) {
	ctx := context.Background()
	defer ctx.Done()

	object, err := service.Client.GetObject(ctx, MinioBucketName, fileUuid, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			log.Println(err)
		}
	}(object)

	objectInfo, err := object.Stat()
	if err != nil {
		return nil, err
	}

	dataBuffer := new(bytes.Buffer)
	copySize, err := dataBuffer.ReadFrom(object)
	if err != nil {
		return nil, err
	}
	if copySize != objectInfo.Size {
		return nil, errors.New("не удалось считать данные файла")
	}

	return dataBuffer.Bytes(), err
}
