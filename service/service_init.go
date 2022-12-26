package service

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/middleware/login"
)

func WebService() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	api := router.Group("/api")
	api.Use(login.IsAuthed())
	api.GET("/ping", Ping)
	api.POST("/login", Login)
	api.GET("/userinfo", UserInfo)

	api.POST("/projects/search", ProjectList)
	api.POST("/projects", CreateProject)
	api.GET("/projects/:projectid", ProjectDetail)

	api.POST("/accounts/search", QueryAccountList)
	api.POST("/accounts", CreateAccount)
	api.DELETE("/accounts", DeleteAccount)
	api.GET("/accounts/info", AccountInfo)
	api.PUT("/accounts/:accountid", UpdateAccount)
	api.PUT("/accounts/:accountid/changepass", ChangeAccountPassword)

	commonApi := api.Group("/common")
	commonApi.GET("/publickey", GetPublicKey)
	commonApi.POST("/register/action", RegisterAction)
	return router
}
