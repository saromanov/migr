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
	"github.com/saromanov/migr/pkg/db"
	"github.com/saromanov/migr/pkg/model"
)

// Run provides starting of migrations
func (a *App) Run(path string) error {
	dirs, err := getMigrationsDirs(path)
	if err != nil {
		return err
	}
	dirs = sortMigrDirs(dirs, 0)

	if err := a.applyMigrations(dirs); err != nil {
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
func sortMigrDirs(dirs []directory, direction uint) []directory {
	sort.Slice(dirs[:], func(i, j int) bool {
		if direction == 0 {
			return dirs[i].timestamp < dirs[j].timestamp
		}
		return dirs[i].timestamp > dirs[j].timestamp
	})
	return dirs
}

// applyMigrations makes migrations
func (a *App) applyMigrations(dirs []directory) error {
	var migrations int
	for i, d := range dirs {
		file, err := ioutil.ReadFile(fmt.Sprintf("./%s/up.sql", d.name))
		if err != nil {
			return errors.Wrap(err, "unable to read up.sql")
		}
		hash, _ := a.hashText(file)
		migr, ok, err := a.isApplyedAndHashed(d, hash)
		if err != nil {
			return errors.Wrap(err, "error on applied migration")
		}
		if ok {
			Info("migration %d is applied", d.timestamp)
			if i+1 < len(dirs) {
				if applied := a.checkNextMigration(dirs[i+1].timestamp); !applied {
					Info("migration %d is not applied. Skip", d.timestamp)
				}
			}
			continue
		}
		if migr != nil {
			if err := a.db.ExecuteCommand(string(file)); err != nil {
				return errors.Wrap(err, fmt.Sprintf("migration %d is not applied", d.timestamp))
			}
			_, err := a.handleApplyingMigration(migr)
			if err != nil {
				return errors.Wrap(err, "unable to update migration")
			}
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

// handleApplyingMigration provides checking of migration
// if migration on Rejected status, then trying to apply this
// an change status to applied
//
// if migration is not exist, then apply this
// if migration on pending status, then skip it
func (a *App) handleApplyingMigration(migr *model.Migration) (bool, error) {
	if migr.Status == db.RejectedStatus {
		return false, a.db.UpdateMigration(migr.ID, true, db.AppliedStatus)
	}
	return true, nil
}

// if migration is applied, then next one should be applied too
// if next one is not applied, then this migration setting
// to the pending status
func (a *App) checkNextMigration(version int64) bool {
	_, ok := a.isApplied(directory{timestamp: version})
	return ok
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

// isApplyedAndHashed checks if migration was already applied and hash sum is equal
func (a *App) isApplyedAndHashed(d directory, hash string) (*model.Migration, bool, error) {
	migr, ok := a.isApplied(d)
	if !ok {
		return migr, false, nil
	}
	if migr != nil && *migr.Hash != hash {
		return nil, false, fmt.Errorf("hash of the migration %d is not equal", d.timestamp)
	}
	return migr, true, nil
}

// isApplyed checks if migration was already applied
func (a *App) isApplied(d directory) (*model.Migration, bool) {
	migr, err := a.db.GetMigrationByTheVersion(d.timestamp)
	if err != nil || migr == nil {
		return nil, false
	}
	if !migr.Applied {
		return migr, false
	}
	return migr, true
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
