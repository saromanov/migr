package app

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Down defines downgrade for migrations
func (a *App) Down(path string) error {
	dirs, err := getMigrationsDirs(path)
	if err != nil {
		return err
	}

	dirs = sortMigrDirs(dirs)

	if err := a.downgradeMigrations(dirs); err != nil {
		return err
	}

	return nil
}

// downgradeMigrations makes migrations
func (a *App) downgradeMigrations(dirs []directory) error {
	for _, d := range dirs {
		file, err := ioutil.ReadFile(fmt.Sprintf("./%s/down.sql", d.name))
		if err != nil {
			return errors.Wrap(err, "unable to read down.sql")
		}

		if err := a.db.ExecuteCommand(string(file)); err != nil {
			return errors.Wrap(err, fmt.Sprintf("migration %d is not applied", d.timestamp))
		}
	}
	return nil
}
