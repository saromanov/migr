package app

import (
	"strings"

	"github.com/pkg/errors"
)

// Down defines downgrade for migrations
func (a *App) Down(path, version string) error {
	dirs, err := getMigrDirs(path)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return errors.New("migr directories is not found")
	}

	var fileName string
	for _, dir := range dirs {
		if strings.Contains(dir.name, version) {
			fileName = dir.name
			break
		}
	}
	if fileName == "" {
		return errors.New("unable to find file")
	}

	return nil
}
