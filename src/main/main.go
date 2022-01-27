package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
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
	//db, err := sql.Open("mysql", "mysql:mysql@tcp(127.0.0.1:3306)/todo")
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

	endpoint := "localhost"
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

	log.Printf("%#v\n", minioClient) // minioClient is now setup

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/getNotes", getNotes)

	err = router.Run(":8222")
	if err != nil {
		return
	}
}
