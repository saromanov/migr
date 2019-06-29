package app_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/app"
)

const (
	basicPath       = "../../testdata/basic"
	migrCreateTable = "CREATE TABLE IF NOT EXISTS migr( id SERIAL, version int8 UNIQUE, changes varchar(128), hash varchar(128), applied bool, error_message varchar(128), failed bool, status varchar(16), created_at int8 )"
)

var (
	appTest *app.App
	db      *sql.DB
)

func init() {
	os.Setenv("MIGR_PATH", basicPath)
	appTest = app.New("postgres", "pinger", "pinger", "pinger", "pinger", 5432)
	dbTmp, err := createTestTable()
	if err != nil {
		fmt.Printf("unable to init db: %v", err)
	}
	db = dbTmp
}

// Create table for tests
func createTestTable() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=migr_test password=migr_test dbname=migr_test")
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err)
	}
	defer db.Close()
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error to ping db")
	}

	return db, nil
}

func createMigrTable(db *sql.DB) error {
	_, err := db.Exec(migrCreateTable)
	if err != nil {
		return fmt.Errorf("unable to create migr table: %v", err)
	}

	return nil
}

func dropMigrTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE migr")
	if err != nil {
		return fmt.Errorf("unable to drop migr table: %v", err)
	}

	return nil
}

func TestRun(t *testing.T) {
	err := appTest.Run(basicPath)
	if err != nil {
		t.Errorf("unable to execute run command: %v", err)
	}
}
