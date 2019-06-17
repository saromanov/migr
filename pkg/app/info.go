package app

import (
	"github.com/pkg/errors"
)

// Info returns information about migrations
func (a *App) Info() error {
	migs, err := a.db.GetMigrationVersions()
	if err != nil {
		return errors.Wrap(err, "unable to get list of migrations")
	}

	for _, mig := range migs {
		Info("migration: %d hash %v applied %v", mig.Version, *mig.Hash, mig.Applied)
	}
	return nil
}
