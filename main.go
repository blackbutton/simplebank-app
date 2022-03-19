package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simplebank-app/api"
	db "simplebank-app/db/sqlc"
	"simplebank-app/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to %s fail", config.DBDriver)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	log.Printf("server listen %s", config.ServerAddress)
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("start server fail: ", err)
	}
}
