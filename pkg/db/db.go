package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
	connString := fmt.Sprintf("user=%s password=%s dbname=%s",
		d.Username, d.Password, d.Database)

	db, err := sql.Open(d.Driver, connString)
	if err != nil {
		return errors.Wrap(err, "error creating connection pool")
	}
	defer db.Close()
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "error to ping db")
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dataBaseTable))
	if err != nil {
		return errors.Wrap(err, "unable to create database")
	}

	_, err = db.Exec(fmt.Sprintf("USE %s", dataBaseTable))
	if err != nil {
		return errors.Wrap(err, "unable to execute")
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s( id integer, version varchar(128), changes varchar(128)", dataBaseTable))
	if err != nil {
		return errors.Wrap(err, "unable to create migr table")
	}
	return nil
}
