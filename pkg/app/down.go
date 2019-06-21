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

	if err := a.downgradeMigrations(sortMigrDirs(dirs, 1)); err != nil {
		return err
	}

	return nil
}

// downgradeMigrations makes migrations
func (a *App) downgradeMigrations(dirs []directory) error {
	for _, d := range dirs {
		if err := a.downgradeMigration(".", d.timestamp); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to downgrade migration %v", d.timestamp))
		}

	}
	return nil
}

// downgradeMigration provides downgrading of migration by the version
func (a *App) downgradeMigration(path string, timestamp int64) error {
	file, err := ioutil.ReadFile(fmt.Sprintf("./%s/down.sql", path))
	if err != nil {
		return errors.Wrap(err, "unable to read down.sql")
	}
	migr, err := a.db.GetMigrationByTheVersion(timestamp)
	if err != nil {
		return errors.Wrap(err, "unable to get migration record")
	}

	if err := a.db.ExecuteCommand(string(file)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("migration %d is not applied", timestamp))
	}

	if err := a.db.WriteMigrationIsApplied(migr.ID, false); err != nil {
		return errors.Wrap(err, fmt.Sprintf("migration %d is not applied", timestamp))
	}
	return nil
}
