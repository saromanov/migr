package app

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Run provides starting of migrations
func (a *App) Run(path string) error {
	dirs, err := getMigrDirs(path)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return errors.New("migr directories is not found")
	}
	dirs = sortMigrDirs(dirs)

	if err := applyMigrations(dirs); err != nil {
		return err
	}
	return nil
}

// getMigrDirs returns dirs which contains
// "migr" on names
func getMigrDirs(path string) ([]directory, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return []directory{}, errors.Wrap(err, "unable to read dir")
	}

	dirs := []directory{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.Contains(name, "migr") {
			parts := strings.Split(name, "_")
			if len(parts) == 0 || len(parts) < 2 {
				continue
			}

			timestamp, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			dirs = append(dirs, directory{
				name:      name,
				timestamp: timestamp,
			})
		}
	}

	return dirs, nil
}

// sortMigrDirs applyins sorting of directories
// by timestamp on the name
func sortMigrDirs(dirs []directory) []directory {
	sort.Slice(dirs[:], func(i, j int) bool {
		return dirs[i].timestamp < dirs[j].timestamp
	})
	return dirs
}

// applyMigrations makes migrations
func applyMigrations(dirs []directory) error {
	for _, d := range dirs {
		file, err := ioutil.ReadFile(fmt.Sprintf("./%s/up.sql", d.name))
		if err != nil {
			return errors.Wrap(err, "unable to read up.sql")
		}
		fmt.Println(string(file))
	}
	return nil
}

// applyCommand provides applying of the sql command
func (a *App) applyCommand(command string) error {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		return errors.Wrap(err, "unable to open connection")
	}

	res, err := db.Exec(command)
	if err != nil {
		return errors.Wrap(err, "unable to execute command")
	}

	return nil
}
