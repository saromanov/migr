package app

import (
	"fmt"

	"github.com/pkg/errors"
)

// DownTo defines downgrade of migration to specific version
func (a *App) DownTo(version string) error {
	dirs, err := getMigrationsDirs(".")
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if fmt.Sprintf("%d", d.timestamp) == version {
			break
		}
		if err := a.downgradeMigration(".", d.timestamp); err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to downgrade migration %v", d.timestamp))
		}
	}

	return nil
}
