package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"todo/src/filestorage"
	"todo/src/controllers"
	"todo/src/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// minioExample()
	var router = note_controller.SetupRouter(note_controller.SetupMiddleware)
	// router = note_controller.SetupMiddleware(router)
	err := router.Run(":8222")
	if err != nil {
		return
	}
}

func minioExample() {
	endpoint := utils.MinioEndpoint
	accessKeyID := utils.MinioAccessKey
	secretAccessKey := utils.MinioSecretKey

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	defer ctx.Done()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		log.Fatalln(err)
		return
	}
	for i := 0; i < len(buckets); i++ {
		log.Println(buckets[i].Name)
	}

	// Upload the zip file
	fileName := "Screenshot_580.png"
	filePath := "C:\\Users\\lappi\\Desktop\\Screenshot_580.png"

	open, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer open.Close()
	stat, err := open.Stat()
	if err != nil {
		log.Fatalln(err)
		return
	}

	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(open)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if size != stat.Size() {
		log.Fatalln(err)
		return
	}

	service := filestorage.MinioServiceImpl{Client: minioClient}
	saveFileUuid, err := service.SaveFile(buf.Bytes(), fileName)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Successfully uploaded %s of size %d\n", fileName, len(buf.Bytes()))

	data, err := service.GetFile(*saveFileUuid)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Successfully get %s of size %d\n", fileName, len(buf.Bytes()))

	name := fileName
	localFile, err := os.Create("C:\\Users\\lappi\\Desktop\\" + *saveFileUuid + name)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err = io.Copy(localFile, bytes.NewReader(data)); err != nil {
		log.Println(err)
		return
	}

	err = service.RemoveFile(*saveFileUuid)
	if err != nil {
		log.Println(err)
		return
	}
}
