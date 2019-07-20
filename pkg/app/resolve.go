package app

import "github.com/saromanov/migr/pkg/model"

// Resolve provides resolving of migrations
func (a *App) Resolve(path string) error {
	dirs, err := getMigrationsDirs(path)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return nil
	}

	dirs = sortMigrDirs(dirs, 0)

	migrations, err := a.db.GetMigrationVersions()
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.Status != "Pending" {
			continue
		}

	}
	return nil
}

func compareDirsAndMigrations(dirs []directory, migrations []*model.Migration) error {
	return nil
}
