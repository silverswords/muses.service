package roommodel

import (
	"database/sql"
	"errors"
	"time"
)

const (
	mysqlCreateRoomTable = iota
	mysqlCreateRoom
	mysqlDeleteRoom
	mysqlGetRooms
	mysqlGetRoomInfo
	mysqlUpdateRoomName
	mysqlUpdateRoomNum
)

var (
	errInvalidMysql = errors.New("affected 0 rows")

	roomsSQLString = []string{
		`CREATE TABLE IF NOT EXISTS rooms (
			id    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name     	VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			num			INT
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO rooms (name, num)  VALUES (?, 0)`,
		`DELETE FROM rooms WHERE id = ? LIMIT 1`,
		`SELECT * FROM rooms`
		`SELECT * FROM rooms WHERE id = ?`
		`Update rooms SET name = ? WHERE id = ?`,
		`Update rooms SET num = ? WHERE id = ?`
	}
)

type Room struct {
	ID uint
	Name string
	Num uint
	CreateTime time.Time
}

func CreateRoomTable(db *sql.DB) error {
	_, err := db.Exec(roomsSQLString[mysqlCreateRoomTable])
	if err != nil {
		return err
	}

	return nil
}

func CreateRoom(db *sql.DB, name string) error {
	_, err := db.Exec(roomsSQLString[mysqlCreateRoom], name)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRoom(db *sql.DB, id int) error {
	_, err := db.Exec(roomsSQLString[mysqlDeleteRoom], id)
	if err != nil {
		return err
	}
	return nil
}

func GetRooms(db *sql.DB) ([]*Room,error) {
	var (
		id uint
		name string
		num uint
		createTime time.Time
		rooms []*Room
	)

	rows, err := db.Query(roomsSQLString[mysqlGetRooms])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &num, &createTime); err != nil {
			return nil, err
		}

		room := &Room{
			ID: id,
			Name:   name,
			Num:       num,
			CreateTime:     createTime,
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func GetRoomInfo(db *sql.DB, id uint) (*Room,error) {
	var (
		id uint
		name string
		num uint
		createTime time.Time
		room *Room
	)

	row, err := db.Query(roomsSQLString[mysqlGetRoomInfo], id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if err := row.Scan(&id, &name, &num, &createTime); err != nil {
		return nil, err
	}

	room := &Room{
		ID: id,
		Name:   name,
		Num:       num,
		CreateTime:     createTime,
	}

	return room, nil
}

func UpdateRoomName(db *sql.DB, id int, name string) error {
	_, err := db.Exec(roomsSQLString[mysqlUpdateRoomName], id, name)
	if err != nil {
		return err
	}

	return nil
}

func JoinRoom(db *sql.DB, id int) error {
	room, err := GetRoomInfo(id)
	if err != nil {
		return err
	}
	
	num := room.Num + 1

	_, err := db.Exec(roomsSQLString[mysqlUpdateRoomNum], num)
	if err != nil {
		return err
	}

	return nil
}

func LeaveRoom(db *sql.DB, id int) error {
	room, err := GetRoomInfo(id)
	if err != nil {
		return err
	}
	
	num := room.Num - 1

	_, err := db.Exec(roomsSQLString[mysqlUpdateRoomNum], num)
	if err != nil {
		return err
	}

	return nil
}
