package main

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/http"
	"os"
	"todo/src/entity"
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
		{Id: 1, NoteGuid: "not guid", Version: 1,
			Text: "text", UserId: 1, Deleted: false, Archive: false,
			NoteFiles: []entity.NoteFile{
				{Id: 1, NoteId: 1, Guid: "note file guid", Filename: "filename"},
			},
		},
	}
	c.IndentedJSON(http.StatusOK, notes)
}

func main() {
	//defer db.Close()
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var version string
	//
	//err2 := db.QueryRow("SELECT VERSION()").Scan(&version)
	//
	//if err2 != nil {
	//	log.Fatal(err2)
	//}
	//
	//fmt.Println("MYSQL VERSION" + version)

	endpoint := "localhost:9000"
	accessKeyID := "minio"
	secretAccessKey := "miniopsw"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := "todo-app-bucket"
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
	uuidWithHyphen := uuid.New().String()
	fileName := "Screenshot_580.png"
	filePath := "C:\\Users\\lappi\\Desktop\\Screenshot_580.png"

	open, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	mtype, err := mimetype.DetectFile(filePath)
	defer open.Close()
	stat, err := open.Stat()
	if err != nil {
		log.Fatalln(err)
		return
	}

	userMetaData := map[string]string{
		"Filename": fileName,
	}
	options := minio.PutObjectOptions{
		ContentType:  mtype.String(),
		UserMetadata: userMetaData,
	}
	info, err := minioClient.PutObject(ctx, bucketName, uuidWithHyphen, open, stat.Size(), options)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Successfully uploaded %s of size %d\n", fileName, info.Size)

	object, err := minioClient.GetObject(context.Background(), bucketName, uuidWithHyphen, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	objectInfo, err := object.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	name := objectInfo.UserMetadata["Filename"]
	localFile, err := os.Create("C:\\Users\\lappi\\Desktop\\" + uuidWithHyphen + name)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = io.Copy(localFile, object); err != nil {
		log.Println(err)
		return
	}

	//router := gin.Default()
	//router.GET("/albums", getAlbums)
	//router.GET("/getNotes", getNotes)
	//
	//err = router.Run(":8222")
	//if err != nil {
	//	return
	//}
}
