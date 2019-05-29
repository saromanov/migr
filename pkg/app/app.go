package app

import (
	"fmt"
	"io/ioutil"
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
	path := fmt.Sprintf("migr_%s_%v", name, time.Now().UnixNano())
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

// Run provides starting of migrations
func (a *App) Run(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return errors.Wrap(err, "unable to read dir")
	}

	for _, f := range files {
		fmt.Println(f)
	}

	return nil
}

func createFile(path string) error {
	if _, err := os.Create(path); err != nil {
		return err
	}

	return nil
}
