package room

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	roommodel "muses.service/model/room"
)

type Controller struct {
	db *sql.DB
}

// New create an external service interface
func NewDB(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	r.POST("/createRoom", c.createRoom)
	r.POST("/deleteRoom", c.deleteRoom)
	// r.POST("/joinRoom", c.joinRoom)
	// r.POST("/leaveRoom", c.leaveRoom)
	r.POST("/updateRoomName", c.updateRoomName)
	r.POST("/getRooms", c.getRooms)
	r.POST("/getRoomByID", c.getRoomInfo)
}

func (c *Controller) createRoom(ctx *gin.Context) {
	var (
		newRoom struct {
			Name string `json:"name"`
		}
	)

	err := ctx.ShouldBind(&newRoom)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = roommodel.CreateRoom(c.db, newRoom.Name)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) deleteRoom(ctx *gin.Context) {
	var (
		newRoom struct {
			ID int `json: "id"`
		}
	)

	err := ctx.ShouldBind(&newRoom)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = roommodel.DeleteRoom(c.db, newRoom.ID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) updateRoomName(ctx *gin.Context) {
	var (
		newRoom struct {
			ID   int    `json: "id"`
			Name string `json: "name"`
		}
	)

	err := ctx.ShouldBind(&newRoom)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = roommodel.UpdateRoomName(c.db, newRoom.ID, newRoom.Name)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) getRooms(ctx *gin.Context) {
	rooms, err := roommodel.GetRooms(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "rooms": rooms})
}

func (c *Controller) getRoomInfo(ctx *gin.Context) {
	var (
		id struct {
			ID uint `json: "id"`
		}
	)

	err := ctx.ShouldBind(&id)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	room, err := roommodel.GetRoomInfo(c.db, id.ID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "room": room})
}

// func (c *Controller) joinRoom(ctx *gin.Context) {

// 	if err := ctx.ShouldBindJSON(&room); err != nil {
// 		ctx.Error(err)
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
// 		return
// 	}

// 	result.Message = "操作成功"
// 	result.Code = http.StatusOK
// 	ctx.JSON(result.Code, gin.H{
// 		"result": result,
// 	})
// }

// func (c *Controller) leaveRoom(ctx *gin.Context) {
// 	room := &Room{}
// 	if err := ctx.ShouldBindJSON(&room); err != nil {
// 		ctx.Error(err)
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
// 		return
// 	}

// 	result.Message = "操作成功"
// 	result.Code = http.StatusOK
// 	ctx.JSON(result.Code, gin.H{
// 		"result": result,
// 	})
// }
