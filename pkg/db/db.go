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
func (d *DB) CreateTable() error {
	if d == nil {
		return errors.New("db is not defined")
	}

	db, err := sql.Open(d.Driver, d.getConnectionString())
	if err != nil {
		return errors.Wrap(err, "error creating connection pool")
	}
	defer db.Close()
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "error to ping db")
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s( id integer, version varchar(128), changes varchar(128))", dataBaseTable))
	if err != nil {
		return errors.Wrap(err, "unable to create migr table")
	}
	return nil
}

// WriteMigrationVersion provides writing of migration version
func (d *DB) WriteMigrationVersion(version string) error {
	connStr := d.getConnectionString()
	db, err := sql.Open(d.Driver, connStr)
	if err != nil {
		return errors.Wrap(err, "unable to open connection")
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO(version, changes) %s($1, $2, $3)", dataBaseTable), version, version)
	if err != nil {
		return fmt.Errorf("unable to execute: %v", err)
	}

	return nil
}

// ExecuteCommand provides execution of the command
func (d *DB) ExecuteCommand(command string) error {
	connString := fmt.Sprintf("user=%s password=%s dbname=%s",
		d.Username, d.Password, d.Database)
	db, err := sql.Open(d.Driver, connString)
	if err != nil {
		return errors.Wrap(err, "unable to open connection")
	}

	_, err = db.Exec(command)
	if err != nil {
		return errors.Wrap(err, "unable to execute command")
	}

	return nil
}

// checkMigrations provides checking of history for migrations
func (d *DB) checkMigrations() error {
	db, err := sql.Open(d.Driver, d.getConnectionString())
	if err != nil {
		return errors.Wrap(err, "error creating connection pool")
	}
	defer db.Close()
	return nil
}

func (d *DB) getConnectionString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s",
		d.Username, d.Password, d.Database)
}
