package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

const dataBaseTable = "migr"

// DB provides handling of db data
type DB struct {
	Username string
	Password string
	Database string
	Driver   string
	Port     int
	Host     string
}

// CreateTable provides creating of the migr table
func CreateTable(d *DB) error {
	if d == nil {
		return errors.New("db is not defined")
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		d.Host, d.Username, d.Password, d.Port, d.Database)

	db, err := sql.Open(d.Driver, connString)
	if err != nil {
		return errors.Wrap(err, "error creating connection pool")
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "error to ping db")
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dataBaseTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf("USE %s", dataBaseTable))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s( id integer, data varchar(32) )", dataBaseTable))
	if err != nil {
		panic(err)
	}
	return nil
}
