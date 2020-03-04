package apis

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"muses.service/handler/student"
	"muses.service/handler/teacher"
	"muses.service/service/eventbus"
)

func InitRouter(bus eventbus.EventBus) *gin.Engine {
	r := gin.Default()
	apiGrp := r.Group("api/v1")

	dbConn, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/test")
	dbmap := &gorp.DbMap{Db: dbConn, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	err = dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	if err != nil {
		panic(err)
	}

	studentConn := student.NewDB(dbmap)
	studentConn.RegisterRouter(apiGrp.Group("/student"))

	teacherConn := teacher.NewDB(dbmap)
	teacherConn.RegisterRouter(apiGrp.Group("/teacher"))

	// roomConn := room.NewDB(dbConn)
	// roomGroup := apiGrp.Group("/room")
	// roomGroup.Use(middleware.MwUser)
	// roomConn.RegisterRouter(roomGroup)

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	return r
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
