package app

import (
	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/model"
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

// GetMigrationsInfo returns list of models for migrations
func (a *App) GetMigrationsInfo() ([]*model.Migration, error) {
	migs, err := a.db.GetMigrationVersions()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of migrations")
	}
	return migs, nil
}
