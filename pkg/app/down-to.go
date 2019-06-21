package app

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// DownTo defines downgrade of migration to specific version
func (a *App) DownTo(version string) error {
	numVersion, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		return errors.Wrap(err, "unable to parse version")
	}
	mig, err := a.db.GetMigrationByTheVersion(numVersion)
	if err != nil {
		return errors.Wrap(err, "unable to find migration")
	}
	if mig == nil {
		return fmt.Errorf("unable to find migration")
	}
	dirs, err := getMigrationsDirs(".")
	if err != nil {
		return err
	}

	dirs = sortMigrDirs(dirs, 1)
	for _, d := range dirs {
		if d.timestamp == mig.Version {
			break
		}
		if err := a.downgradeMigration(d.name, d.timestamp); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to downgrade migration %v", d.timestamp))
		}
		Info("migration %d is downgraded", d.timestamp)
	}

	return nil
}
