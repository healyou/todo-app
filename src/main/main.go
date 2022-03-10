package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"todo/src/entity"
	"todo/src/filestorage"
	"todo/src/controllers"
	"todo/src/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//go get github.com/minio/minio-go/v7@v7.0.21

// album represents data about a record album.
//type album struct {
//	ID     string  `json:"id"`
//	Title  string  `json:"title"`
//	Artist string  `json:"artist"`
//	Price  float64 `json:"price"`
//}

// albums slice to seed record album data.
var albums = []entity.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getNotes(c *gin.Context) {
	var notes = []entity.Note{
		//{Id: 1, NoteGuid: "not guid", Version: 1,
		//	Text: "text", UserId: 1, Deleted: false, Archive: false,
		//	NoteFiles: []entity.NoteFile{
		//		{Id: 1, NoteId: 1, Guid: "note file guid", Filename: "filename"},
		//	},
		//},
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func main() {
	// minioExample()
	router := gin.Default()
	router.POST("/notes/getActualNote", note_controller.GetActualNote)
	router.POST("/notes/saveNote", note_controller.SaveNote)
	router.POST("/notes/getUserNotes", note_controller.GetUserNotes)
	router.POST("/notes/downNoteVersion", note_controller.DownNoteVersion)
	router.POST("/notes/upNoteVersion", note_controller.UpNoteVersion)

	router.GET("/albums", getAlbums)
	router.GET("/getNotes", getNotes)

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
