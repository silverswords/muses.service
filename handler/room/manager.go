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
	Rooms            map[string]room
	connetionManager connection.Manager
}

type room struct {
	RoomID  string
	Persons []string
}

// NewManger -
func NewManger(connetionManager *connection.Manager) *Manager {
	return &Manager{
		Rooms:            make(map[string]room),
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

	fmt.Println(roomID.ID)

	m.Rooms[roomID.ID] = room{
		RoomID: roomID.ID,
	}

	fmt.Println(m.Rooms)
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

	fmt.Println(m.Rooms)
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

	_ = append(room.Persons, param.UserID)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (m *Manager) sendMes(ctx *gin.Context) {
	var (
		msg struct {
			userid  string `json: "id"`
			roomid  string `json: "roomid"`
			content string `json: "content"`
		}
	)

	err := ctx.ShouldBind(&msg)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	room, ok := m.Rooms[msg.roomid]
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "room not exited",
		})
		ctx.Error(err)
		return
	}

	ret := make([]string, 0, len(room.Persons))
	for _, val := range room.Persons {
		if val != msg.userid {
			ret = append(ret, val)
		}
	}

	go m.connetionManager.SubMsg(ret, msg.content)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
