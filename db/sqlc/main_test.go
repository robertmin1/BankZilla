package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

)

var testQueries *Queries
var testDB *sql.DB

const (
	driverName = "postgres"
)

func TestMain(m *testing.M)  {
	var err error

	// Read the POSTGRES_PASSWORD environment variable
    	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

    	// Construct the dataSourceName
   	dataSourceName := "postgresql://root:" + postgresPassword + "@localhost:5432/simple_bank?sslmode=disable"

	testDB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
