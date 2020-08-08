package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBTools interface {
	CheckDB(retry int) error
	CreateExporter() error
}

func NewDBTools(db *sql.DB) DBTools {
	return &DBToolsType{DB: db}
}

type DBToolsType struct {
	DB *sql.DB
}

func (db *DBToolsType) CheckDB(retry int) error {
	for i := 0; i < retry; i++ {
		err := db.DB.Ping()
		if err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return errors.New("connection failed")
}

func (db *DBToolsType) CreateExporter() error {
	sqls := []string{
		"create user if not exists 'exporter'@'localhost' identified by 'exporter' WITH MAX_USER_CONNECTIONS 3",
		"create user if not exists 'exporter'@'::1' identified by 'exporter' WITH MAX_USER_CONNECTIONS 3",
		"GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'localhost'",
		"GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'::1'",
	}
	for _, v := range sqls {
		_, err := db.DB.Exec(v)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return err
		}
	}
	return nil
}
