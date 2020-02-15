package apis

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"muses.service/handler"
	"muses.service/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	dbConn, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}

	adminConn := handler.NewUserDB(dbConn)
	adminConn.RegisterRouter(r.Group("/api/v1"))

	apiGrp(r)
	return r
}

func apiGrp(r *gin.Engine) {
	apiGrp := r.Group("/api/v1")

	// users
	userGrp := apiGrp.Group("/user")
	userGrp.Use(middleware.MwUser)
	{
		userGrp.POST("/joinRoom", handler.Joinroom)
		userGrp.POST("/leaveRoom")
		userGrp.POST("/sendMsg")

		userGrp.POST("/removeUser")
	}

	// room
	// roomGrp := apiGrp.Group("/room")
	// {
	// 	roomGrp.POST("/getRoomList")
	// 	roomGrp.POST("/getRoomInfo")

	// 	roomGrp.POST("/createRoom")
	// 	roomGrp.POST("/removeRoom")
	// }

	// // file
	// fileGrp := apiGrp.Group("/file")
	// {
	// 	fileGrp.POST("/listFile")
	// 	fileGrp.POST("/uploadFile")
	// 	fileGrp.POST("/removeFile")
	// }
}
