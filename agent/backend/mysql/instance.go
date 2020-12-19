package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// DefaultDelay is the default poll interval to check an Instance
const DefaultDelay = 5 * time.Second

// NewInstance provides the interfaces required to start an instance
func NewInstance(db *sql.DB) *Instance {
	return &Instance{DB: db}
}

// Instance is the interface implementation for MySQL
type Instance struct {
	DB *sql.DB
}

// Check wait for the instance to start
func (db *Instance) Check(retry int) error {
	for i := 0; i < retry; i++ {
		err := db.DB.Ping()
		if err == nil {
			return nil
		}
		time.Sleep(DefaultDelay)
	}
	return errors.New("connection failed")
}

// Initialize creates a set of queries to connect to the database from tools
func (db *Instance) Initialize() error {
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
