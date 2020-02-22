package room

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// Manager -
type RoomManager interface {
	InitRoom(name string) error
	RemoveRoom(roomID string) error
	ListRooms() map[string]Room
	JoinRoom(roomID string, uid string) error
	LeaveRoom(roomID string, uid string) error
	SendMessage(roomID string, msg string, uid string) error
}

// Manager -
type Manager struct {
	rooms map[string]Room
}

// Room -
type Room struct {
	id    string
	conns []string
	name  string
	num   int
}

// Init -
func Init() *Manager {
	return &Manager{
		rooms: make(map[string]Room),
	}
}

// InitRoom -
func (manager *Manager) InitRoom(name string) error {
	for _, v := range manager.rooms {
		if name == v.name {
			return errors.New("this room is already exist")
		}
	}

	newID := uuid.NewV4()
	manager.rooms[newID.String()] = Room{
		id:    newID.String(),
		conns: make([]string, 1),
		name:  name,
		num:   0,
	}

	return nil
}

// RemoveRoom -
func (manager *Manager) RemoveRoom(roomID string) error {
	_, ok := manager.rooms[roomID]
	if !ok {
		return errors.New("room not exsit")
	}

	delete(manager.rooms, roomID)
	return nil
}

// ListRooms -
func (manager *Manager) ListRooms() map[string]Room {
	return manager.rooms
}

// JoinRoom -
func (manager *Manager) JoinRoom(roomID string, uid string) error {
	room, ok := manager.rooms[roomID]
	if !ok {
		return errors.New("room not exsit")
	}

	for _, conn := range room.conns {
		if conn == uid {
			return errors.New("already in room")
		}
	}

	room.conns = append(room.conns, uid)
	return nil
}

// LeaveRoom -
func (manager *Manager) LeaveRoom(roomID string, uid string) error {
	room, ok := manager.rooms[roomID]
	if !ok {
		return errors.New("room not exsit")
	}

	for i, conn := range room.conns {
		if conn == uid {
			room.conns = append(room.conns[:i], room.conns[i+1:]...)
			return nil
		}
	}

	return errors.New("person not exsit")
}

// SendMessage -
func (manager *Manager) SendMessage(roomID string, msg string, uid string) error {
	room, ok := manager.rooms[roomID]
	if !ok {
		return errors.New("room not exsit")
	}

	conns := make([]string, 1)
	for i, conn := range room.conns {
		if conn == uid {
			conns = append(room.conns[:i], room.conns[i+1:]...)
		}
	}

	// send to connctionManager
	// eventBus.send("sendMsg", conns)
	return nil
}
