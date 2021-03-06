package main

import (
	"database/sql"
	"log"

	"github.com/blaqkube/mysql-operator/agent/backend"
	"github.com/blaqkube/mysql-operator/agent/backend/blackhole"
	"github.com/blaqkube/mysql-operator/agent/backend/gcp"
	"github.com/blaqkube/mysql-operator/agent/backend/mysql"
	"github.com/blaqkube/mysql-operator/agent/backend/s3"
	"github.com/blaqkube/mysql-operator/agent/cmd"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/")
	if err != nil {
		log.Printf("Error opening database: %v", err)
	}
	instance := mysql.NewInstance(db)
	backup := mysql.NewBackup()

	storages := map[string]backend.Storage{
		"s3":        s3.NewStorage(),
		"blackhole": blackhole.NewStorage(),
		"gcp":       gcp.NewStorage(),
	}

	cmd.Execute(backup, db, instance, storages)
}
