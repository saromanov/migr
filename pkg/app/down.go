package app

import (
	"github.com/pkg/errors"
)

// Down defines downgrade for migrations
func (a *App) Down(path string) error {
	dirs, err := getMigrDirs(path)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return errors.New("migr directories is not found")
	}

	return nil
}
