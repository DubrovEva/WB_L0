package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

type Database struct {
	Conn  *sql.DB
	Cache *sync.Map
}

func InitializeDB(username, password, dbname, host string, port int64) (*Database, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	var cache sync.Map
	db := Database{Conn: conn, Cache: &cache}

	err = db.Conn.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Database connection established")

	err = db.UpdateCache()
	if err != nil {
		log.Printf("Can't update cache")
	} else {
		log.Printf("Update cache")
	}

	return &db, nil
}
