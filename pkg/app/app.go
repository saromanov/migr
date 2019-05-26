package app

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

// App provides impementation of the main logic
// for the app
type App struct {
	driver string
}

// New creates app
func New(driver string) *App {
	return &App{
		driver: driver,
	}
}

// Create provides creating of the two files
// for migrations up and down. Also, it creates
// records on db
func (a *App) Create(name string) error {
	if err := os.Mkdir(fmt.Sprintf("%s_%v", name, time.Now().UnixNano()), 0755); err != nil {
		return errors.Wrap(err, "unable to create dir")
	}

	if err := createFile(fmt.Sprintf("%s/up.sql")); err != nil {
		return errors.Wrap(err, "unable to create up.sql")
	}

	if err := createFile(fmt.Sprintf("%s/down.sql")); err != nil {
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
