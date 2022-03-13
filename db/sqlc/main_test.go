package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testStore *Store

func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("conect database fail with error %v", err)
	}
	testStore = NewStore(conn)
	testQueries = testStore.Queries
	os.Exit(m.Run())
}
