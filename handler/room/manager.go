package room

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"muses.service/service/connection"
)

// Manager -
type Manager struct {
	Rooms map[string]room
}

type room struct {
	RoomID  string
	Persons []string
}

// NewManger -
func NewManger() *Manager {
	return &Manager{
		Rooms: make(map[string]room),
	}
}

// RegisterRouter -
func (m *Manager) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	NewManger()

	r.POST("/openroom", m.openRoom)
	r.POST("/joinroom", m.joinRoom)
}

// openRoom -
func (m *Manager) openRoom(ctx *gin.Context) {
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

	m.Rooms[roomID.ID] = room{
		RoomID: roomID.ID,
	}
}

func (m *Manager) joinRoom(ctx *gin.Context) {
	var (
		param struct {
			UserID string `json:"id"`
			RoomID string `json:"roomid"`
		}
	)

	err := ctx.ShouldBind(&param)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	room, ok := m.Rooms[param.RoomID]
	if !ok {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "room not exited",
		})
		return
	}

	_, ok = connection.Manager.Connections[param.UserID]
	if !ok {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "connection not exited",
		})
		return
	}

	_ = append(room.Persons, param.UserID)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
