package usermodel

import (
	"database/sql"
	"errors"

	"github.com/fengyfei/comet/pkgs/salt"
)

const (
	mysqlUserCreateTable = iota
	mysqlUserInsert
	mysqlUserLogin
)

var (
	errInvalidMysql = errors.New("affected 0 rows")
	errLoginFailed  = errors.New("invalid username or password")

	adminSQLString = []string{
		`CREATE TABLE IF NOT EXISTS user (
			id    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
			name     	VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			password 	VARCHAR(512) NOT NULL DEFAULT ' ',
			created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO user (name,password)  VALUES (?,?)`,
		`SELECT id,password FROM user WHERE name = ? LOCK IN SHARE MODE`,
	}
)

// CreateTable create admin table.
func CreateTable(db *sql.DB, name, password *string) error {
	_, err := db.Exec(adminSQLString[mysqlUserCreateTable])
	if err != nil {
		return err
	}
	return nil
}

//Create create an administrative user
func Create(db *sql.DB, name, password *string) error {
	hash, err := salt.Generate(password)
	if err != nil {
		return err
	}

	result, err := db.Exec(adminSQLString[mysqlUserInsert], name, hash)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

//Login the administrative user logins
func Login(db *sql.DB, name, password *string) (uint32, error) {
	var (
		id  uint32
		pwd string
	)

	err := db.QueryRow(adminSQLString[mysqlUserLogin], name).Scan(&id, &pwd)
	if err != nil {
		return 0, err
	}

	if !salt.Compare([]byte(pwd), password) {
		return 0, errLoginFailed
	}

	return id, nil
}
