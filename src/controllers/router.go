package controllers

import (
	"todo/src/controllers/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	/* Определение привилегий пользователя */
	router.Use(middleware.AuthorizationMiddleware)

	router.POST("/notes/getActualNote", GetActualNote)
	router.POST("/notes/saveNote", SaveNote)
	router.POST("/notes/getUserNotes", GetUserNotes)
	router.POST("/notes/downNoteVersion", DownNoteVersion)
	router.POST("/notes/upNoteVersion", UpNoteVersion)
	return router
}