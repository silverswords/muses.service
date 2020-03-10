package room

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gorp.v1"
)

// Controller -
type Controller struct {
	dbmap *gorp.DbMap
}

// Room -
type Room struct {
	RoomID    string
	Name      string
	Info      string
	Location  string
	MaxNumber int64
	Created   int64
}

// NewDB -
func NewDB(dbmap *gorp.DbMap) *Controller {
	return &Controller{
		dbmap: dbmap,
	}
}

// RegisterRouter -
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	c.dbmap.AddTableWithName(Room{}, "room").SetKeys(false, "RoomID")

	r.POST("/createRoom", c.createRoom)
	r.POST("/removeRoom", c.removeRoom)
	r.GET("/listRoom", c.listRoom)
	r.POST("/modifyRoom", c.modifyRoom)
}

func (c *Controller) createRoom(ctx *gin.Context) {
	var (
		roomBasic struct {
			Name      string `json:"name"`
			Info      string `json:"info"`
			Location  string `json:"loacton"`
			MaxNumber int64  `json:"maxNumber"`
		}
	)

	err := ctx.ShouldBind(&roomBasic)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	room := Room{
		RoomID:    uuid.NewV4().String(),
		Name:      roomBasic.Name,
		Info:      roomBasic.Info,
		Location:  roomBasic.Location,
		MaxNumber: roomBasic.MaxNumber,
		Created:   time.Now().UnixNano(),
	}

	err = c.dbmap.Insert(&room)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (c *Controller) removeRoom(ctx *gin.Context) {
	var (
		roomID struct {
			ID string `json:"id"`
		}
	)

	err := ctx.ShouldBind(&roomID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	room := Room{
		RoomID: roomID.ID,
	}

	_, err = c.dbmap.Delete(&room)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (c *Controller) listRoom(ctx *gin.Context) {
	var rooms []Room
	_, err := c.dbmap.Select(&rooms, "select * from room")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   rooms,
	})
}

func (c *Controller) modifyRoom(ctx *gin.Context) {
	var (
		roomBasic struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Info      string `json:"info"`
			Location  string `json:"loaction"`
			MaxNumber int64  `json:"maxNumber"`
		}
	)

	err := ctx.ShouldBind(&roomBasic)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	room := Room{
		RoomID:    roomBasic.ID,
		Name:      roomBasic.Name,
		Info:      roomBasic.Info,
		Location:  roomBasic.Location,
		MaxNumber: roomBasic.MaxNumber,
	}

	_, err = c.dbmap.Update(room)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
