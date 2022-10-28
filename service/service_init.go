package service

import (
	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/middleware/login"
)

func WebService() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	api.Use(login.IsAuthed())
	api.GET("/ping", Ping)
	api.POST("/login", Login)
	api.GET("/userinfo", UserInfo)

	return router
}
