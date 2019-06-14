package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/model"
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

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s( id SERIAL, version int8 UNIQUE, changes varchar(128), hash varchar(128), applied bool, error_message varchar(128), failed bool )", dataBaseTable))
	if err != nil {
		return errors.Wrap(err, "unable to create migr table")
	}
	return nil
}

// CreateMigrationVersion provides beginning of the phase for create migration version
func (d *DB) CreateMigrationVersion(version string) (int64, error) {
	connStr := d.getConnectionString()
	db, err := sql.Open(d.Driver, connStr)
	if err != nil {
		return 0, errors.Wrap(err, "unable to open connection")
	}

	if err := d.CreateTable(); err != nil {
		return 0, errors.Wrap(err, "unable to create table")
	}

	var id int64
	err = db.QueryRow(fmt.Sprintf("INSERT INTO %s(version, changes, applied, failed) VALUES ($1, $2, $3, $4) RETURNING id", dataBaseTable), version, version, false, false).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "unable to insert data")
	}
	return id, nil
}

// WriteMigrationVersion provides writing of migration version
func (d *DB) WriteMigrationVersion(id int64) error {
	if id == 0 {
		return fmt.Errorf("id is not defined")
	}
	connStr := d.getConnectionString()
	db, err := sql.Open(d.Driver, connStr)
	if err != nil {
		return errors.Wrap(err, "unable to open connection")
	}

	_, err = db.Exec("UPDATE migr SET applied = $1 WHERE id = $2", true, id)
	if err != nil {
		return fmt.Errorf("unable to execute: %v", err)
	}

	return nil
}

// GetMigrationVersions returns list of migrations
func (d *DB) GetMigrationVersions() ([]*model.Migration, error) {
	connStr := d.getConnectionString()
	db, err := sql.Open(d.Driver, connStr)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open connection")
	}

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", dataBaseTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	migs := make([]*model.Migration, 0)
	for rows.Next() {
		mig := new(model.Migration)
		err := rows.Scan(&mig.ID, &mig.Version, &mig.Changes)
		if err != nil {
			return nil, err
		}
		migs = append(migs, mig)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return migs, nil
}

// GetMigrationByTheVersion returns migration by the version
func (d *DB) GetMigrationByTheVersion(version int64) (*model.Migration, error) {
	connStr := d.getConnectionString()
	db, err := sql.Open(d.Driver, connStr)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open connection")
	}

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE version = $1", dataBaseTable), version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mig := new(model.Migration)
	for rows.Next() {
		err := rows.Scan(&mig.ID, &mig.Version, &mig.Changes, &mig.Hash, &mig.Applied, &mig.ErrorMessage, &mig.Failed)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if mig.ID == 0 {
		return nil, errors.Wrap(err, "unable to find migration")
	}
	return mig, nil
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
