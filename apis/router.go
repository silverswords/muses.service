package apis

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"muses.service/handler/room"
	"muses.service/handler/user"
	"muses.service/middleware"
	"muses.service/service/eventbus"
)

func InitRouter(bus eventbus.EventBus) *gin.Engine {
	r := gin.Default()
	apiGrp := r.Group("api/v1")

	dbConn, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}

	userConn := user.NewDB(dbConn)
	userConn.RegisterRouter(apiGrp.Group("/user"))

	roomConn := room.NewDB(dbConn)
	roomGroup := apiGrp.Group("/room")
	roomGroup.Use(middleware.MwUser)
	roomConn.RegisterRouter(roomGroup)

	return r
}
