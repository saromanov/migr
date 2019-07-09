package app

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/db"
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
		if err := a.downgradeMigration(d.name, d.timestamp); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to downgrade migration %v", d.timestamp))
		}
		Info("migration %d is downgraded", d.timestamp)
	}
	return nil
}

// downgradeMigration provides downgrading of migration by the version
func (a *App) downgradeMigration(name string, timestamp int64) error {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/down.sql", a.path, name))
	if err != nil {
		return errors.Wrap(err, "unable to read down.sql")
	}
	migr, err := a.db.GetMigrationByTheVersion(timestamp)
	if err != nil {
		return errors.Wrap(err, "unable to get migration record")
	}
	if !migr.Applied {
		return nil
	}
	if err := a.db.ExecuteCommand(string(file)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("migration %d is not applied", timestamp))
	}

	if err := a.db.UpdateMigration(migr.ID, false, db.RejectedStatus); err != nil {
		return errors.Wrap(err, fmt.Sprintf("migration %d is not applied", timestamp))
	}
	return nil
}
