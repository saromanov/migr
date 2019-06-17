package app

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/model"
)

// Run provides starting of migrations
func (a *App) Run(path string) error {
	dirs, err := getMigrationsDirs(path)
	if err != nil {
		return err
	}
	dirs = sortMigrDirs(dirs)

	if err := a.applyMigrations(dirs, "up.sql"); err != nil {
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
		if !strings.Contains(name, "migr") {
			continue
		}
		parts := strings.Split(name, "_")
		if len(parts) == 0 || len(parts) < 2 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			continue
		}
		dirs = append(dirs, directory{
			name:      name,
			timestamp: timestamp,
		})
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
func (a *App) applyMigrations(dirs []directory, fileName string) error {
	var migrations int
	for _, d := range dirs {
		file, err := ioutil.ReadFile(fmt.Sprintf("./%s/%s", d.name, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to read %s", fileName))
		}

		hash, _ := a.hashText(file)
		ok, err := a.isApplyed(d, hash)
		if err != nil {
			return errors.Wrap(err, "error on applied migration")
		}
		if ok {
			Info("migration %d is applied", d.timestamp)
			continue
		}
		Info("trying to apply migration %d", d.timestamp)

		id, err := a.db.CreateMigrationVersion(fmt.Sprintf("%d", d.timestamp))
		if err != nil {
			return errors.Wrap(err, "unable to create migration record")
		}

		if err := a.db.ExecuteCommand(string(file)); err != nil {
			return errors.Wrap(err, fmt.Sprintf("migration %d is not applied", d.timestamp))
		}

		err = a.db.WriteMigrationVersion(id, hash)
		if err != nil {
			return errors.Wrap(err, "unable to create migration record")
		}

		Info("migration %d is applied", d.timestamp)
		migrations++
	}
	if migrations > 0 {
		Info("%d migrations is applied", migrations)
	}
	return nil
}

// applyMigration provides applying of migration
func (a *App) applyMigration(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "unable to read up.sql")
	}

	if err := a.db.ExecuteCommand(string(file)); err != nil {
		return errors.Wrap(err, "migrations is not applied")
	}

	return nil
}

// get list of applied migrations from db
func (a *App) getAppliedMigrations(dbname string) ([]*model.Migration, error) {
	migs, err := a.db.GetMigrationVersions()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get list of migrations")
	}
	return migs, nil
}

// isApplyed checks if migration was already applied
func (a *App) isApplyed(d directory, hash string) (bool, error) {
	migr, err := a.db.GetMigrationByTheVersion(d.timestamp)
	if err != nil {
		return false, nil
	}
	if migr == nil {
		return false, nil
	}

	if *migr.Hash != hash {
		return false, fmt.Errorf("hash of the migration %d is not equal", d.timestamp)
	}
	return true, nil
}

// hashText returns hash of the up file
func (a *App) hashText(text []byte) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write(text)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
}
