package apis

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"muses.service/handler/room"
	"muses.service/handler/student"
	"muses.service/handler/teacher"
	"muses.service/service/connection"
)

// InitRouter -
func InitRouter() *gin.Engine {
	r := gin.Default()
	apiGrp := r.Group("api/v1")

	dbConn, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/test")
	dbmap := &gorp.DbMap{Db: dbConn, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	err = dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	if err != nil {
		panic(err)
	}

	// connection
	connetionManager := connection.NewConnectionManager()
	go connetionManager.Run()

	// update WS
	r.GET("/ws", connetionManager.UpGraderWs)

	// student apis
	studentConn := student.NewDB(dbmap)
	studentConn.RegisterRouter(apiGrp.Group("/student"))

	// teacher apis
	teacherConn := teacher.NewDB(dbmap)
	teacherConn.RegisterRouter(apiGrp.Group("/teacher"))

	// roomManager apis
	roomManger := room.NewManger(connetionManager)
	roomManger.RegisterRouter(apiGrp.Group("/room"))

	// room apis
	roomConn := room.NewDB(dbmap)
	roomConn.RegisterRouter(apiGrp.Group("/room"))

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	return r
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
