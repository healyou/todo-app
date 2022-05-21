package controllers

import (
	"todo/src/controllers/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	mainRouterGroup := router.Group("/notes-api")

	/* Определение привилегий пользователя */
	mainRouterGroup.Use(middleware.AuthorizationMiddleware)

	mainRouterGroup.POST("/notes/getActualNote", GetActualNote)
	mainRouterGroup.POST("/notes/saveNote", SaveNote)
	mainRouterGroup.POST("/notes/getNoteFileBody", DownloadNoteFile)
	mainRouterGroup.POST("/notes/getUserNotes", GetUserNotes)
	mainRouterGroup.POST("/notes/getLastUserNoteMainInfo", GetLastUserNoteMainInfo)
	mainRouterGroup.POST("/notes/downNoteVersion", DownNoteVersion)
	mainRouterGroup.POST("/notes/upNoteVersion", UpNoteVersion)
	mainRouterGroup.POST("/notes/getNoteVersionHistory", GetNoteVersionHistory)
	return router
}