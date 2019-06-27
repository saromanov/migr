package app

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/db"
)

// App provides impementation of the main logic
// for the app
type App struct {
	driver string
	db     *db.DB
	path   string
}

type directory struct {
	name      string
	timestamp int64
}

// New creates app
func New(driver, username, password, dbname, host string, port int) *App {
	path := os.Getenv("MIGR_PATH")
	if path == "" {
		path = "."
	}
	return &App{
		driver: driver,
		path:   path,
		db: &db.DB{
			Username: username,
			Password: password,
			Database: dbname,
			Host:     host,
			Port:     port,
			Driver:   driver,
		},
	}
}

// Create provides creating of the two files
// for migrations up and down. Also, it creates
// records on db
func (a *App) Create(name string) error {
	path := fmt.Sprintf("migr_%s_%v", name, time.Now().UnixNano())
	if err := a.db.CreateTable(); err != nil {
		return errors.Wrap(err, "unable to create migr table")
	}
	if err := os.Mkdir(path, 0755); err != nil {
		return errors.Wrap(err, "unable to create dir")
	}

	if err := createFile(fmt.Sprintf("%s/up.sql", path)); err != nil {
		return errors.Wrap(err, "unable to create up.sql")
	}

	if err := createFile(fmt.Sprintf("%s/down.sql", path)); err != nil {
		return errors.Wrap(err, "unable to create down.sql")
	}
	return nil
}

func createFile(path string) error {
	if _, err := os.Create(path); err != nil {
		return err
	}

	return nil
}

func getMigrationsDirs(path string) ([]directory, error) {
	dirs, err := getMigrDirs(path)
	if err != nil {
		return []directory{}, err
	}
	if len(dirs) == 0 {
		return []directory{}, errors.New("migr directories is not found")
	}

	return dirs, nil
}
