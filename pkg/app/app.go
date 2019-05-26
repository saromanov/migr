package app

import (
	"fmt"
	"os"
	"time"
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
	if err := os.Mkdir(fmt.Sprintf("%s_%v", name, time.Now().UnixNano()), os.Perm); err != nil {
		return err
	}

	if err := createFile(fmt.Sprintf("%s/up.sql")); err != nil {
		return err
	}

	if err := createFile(fmt.Sprintf("%s/down.sql")); err != nil {
		return err
	}
	return nil
}

func createFile(path string) error {
	if _, err := os.Create(path); err != nil {
		return err
	}

	return nil
}
