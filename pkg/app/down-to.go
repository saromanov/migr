package app

import "fmt"

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

	}

	return nil
}
