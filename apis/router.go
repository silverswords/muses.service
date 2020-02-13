package apis

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	apiGrp(r)

	return r
}

func apiGrp(r *gin.Engine) {
	apiGrp := r.Group("/api/v1")
	apiGrp.POST("/login")
	apiGrp.POST("/register")
	apiGrp.POST("/default") // 游客模式

	// users
	userGrp := apiGrp.Group("/user")
	// userGrp.Use(mdw.AuthJWT())
	{
		userGrp.POST("/joinRoom")
		userGrp.POST("/leaveRoom")
		userGrp.POST("/sendMsg")

		userGrp.POST("/removeUser")
	}

	// room
	roomGrp := apiGrp.Group("/room")
	{
		roomGrp.POST("/getRoomList")
		roomGrp.POST("/getRoomInfo")

		roomGrp.POST("/createRoom")
		roomGrp.POST("/removeRoom")
	}

	// file
	fileGrp := apiGrp.Group("/file")
	{
		fileGrp.POST("/listFile")
		fileGrp.POST("/uploadFile")
		fileGrp.POST("/removeFile")
	}
}
