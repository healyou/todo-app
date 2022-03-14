package note_controller

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/notes/getActualNote", GetActualNote)
	router.POST("/notes/saveNote", SaveNote)
	router.POST("/notes/getUserNotes", GetUserNotes)
	router.POST("/notes/downNoteVersion", DownNoteVersion)
	router.POST("/notes/upNoteVersion", UpNoteVersion)
	return router
}