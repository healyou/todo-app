package filestorage

type MinioService interface {
	SaveFile(data []byte, filename string) (*string, error)

	GetFile(fileUuid string) ([]byte, error)
}
