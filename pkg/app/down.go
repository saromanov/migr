package app

// Down defines downgrade for migrations
func (a *App) Down(path string) error {
	dirs, err := getMigrationsDirs(path)
	if err != nil {
		return err
	}

	dirs = sortMigrDirs(dirs)

	if err := a.applyMigrations(dirs, "down.sql"); err != nil {
		return err
	}

	return nil
}
