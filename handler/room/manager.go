package room

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"muses.service/service/connection"
)

// Manager -
type Manager struct {
	Rooms            map[string]Room
	connetionManager connection.Manager
}

type Room struct {
	RoomID  string
	Persons []string
}

// NewManger -
func NewManger(connetionManager *connection.Manager) *Manager {
	return &Manager{
		Rooms:            make(map[string]Room),
		connetionManager: *connetionManager,
	}
}

// RegisterRouter -
func (m *Manager) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	r.POST("/openroom", m.openRoom)
	r.POST("/joinroom", m.joinRoom)
	r.POST("/sendmsg", m.sendMes)
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

	m.Rooms[roomID.ID] = Room{
		RoomID:  roomID.ID,
		Persons: make([]string, 0),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "room not exited",
		})
		ctx.Error(err)
		return
	}

	// _, ok = connection.Manager.Connections[param.UserID]
	// if !ok {
	// 	ctx.Error(err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"status": http.StatusBadRequest,
	// 		"msg":    "connection not exited",
	// 	})
	// 	return
	// }
	m.Rooms[param.RoomID].Persons = append(m.Rooms[param.RoomID].Persons, param.UserID)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (m *Manager) sendMes(ctx *gin.Context) {
	var (
		msg struct {
			ID      string `json: "id"`
			RoomID  string `json: "roomid"`
			Content string `json: "content"`
		}
	)

	err := ctx.ShouldBind(&msg)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	fmt.Println(msg.RoomID)
	room, ok := m.Rooms[msg.RoomID]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "room not exited",
		})
		ctx.Error(err)
		return
	}

	go m.connetionManager.SubMsg(room.Persons, msg.Content)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
