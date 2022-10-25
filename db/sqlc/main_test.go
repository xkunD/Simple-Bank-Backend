package db_test

import (
	"database/sql"
	db "go-simple-bank/db/sqlc"
	"go-simple-bank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = db.New(testDb)

	os.Exit(m.Run())
}
