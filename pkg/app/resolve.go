package app

// Resolve provides resolving of migrations
func (a *App) Resolve(path string) error {
	dirs, err := getMigrationsDirs(path)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return nil
	}
	return nil
}
