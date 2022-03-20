package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"simplebank-app/util"
	"testing"
)

var testQueries Querier
var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config file: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("conect database fail with error %v", err)
	}
	testStore = NewStore(conn)
	testQueries = testStore
	os.Exit(m.Run())
}
