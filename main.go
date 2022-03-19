package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simplebank-app/api"
	db "simplebank-app/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable"
	address  = ":9000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to %s fail", dbDriver)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	log.Printf("server listen %s", address)
	if err := server.Start(address); err != nil {
		log.Fatal("start server fail: ", err)
	}
}
