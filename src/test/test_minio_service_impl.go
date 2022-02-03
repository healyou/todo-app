package test

import (
	"github.com/google/uuid"
)

type MinioServiceImplTest struct {
}

func (service MinioServiceImplTest) SaveFile(data []byte, filename string) (*string, error) {
	uuidWithHyphen := uuid.New().String()
	return &uuidWithHyphen, nil
}

func (service MinioServiceImplTest) GetFile(fileUuid string) ([]byte, error) {
	return []byte{}, nil
}
