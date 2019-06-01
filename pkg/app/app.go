package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/db"
)

// App provides impementation of the main logic
// for the app
type App struct {
	driver   string
	password string
	username string
	dbname   string
	port     int
	host     string
}

type directory struct {
	name      string
	timestamp int64
}

// New creates app
func New(driver, username, password, dbname, host string, port int) *App {
	return &App{
		driver:   driver,
		username: username,
		password: password,
		dbname:   dbname,
		host:     host,
		port:     port,
	}
}

// Create provides creating of the two files
// for migrations up and down. Also, it creates
// records on db
func (a *App) Create(name string) error {
	path := fmt.Sprintf("migr_%s_%v", name, time.Now().UnixNano())
	if err := db.CreateTable(&db.DB{
		Username: a.username,
		Password: a.password,
		Database: a.dbname,
		Host:     a.host,
		Port:     a.port,
		Driver:   a.driver,
	}); err != nil {
		return errors.Wrap(err, "unable to create migr table")
	}
	if err := os.Mkdir(path, 0755); err != nil {
		return errors.Wrap(err, "unable to create dir")
	}

	if err := createFile(fmt.Sprintf("%s/up.sql", path)); err != nil {
		return errors.Wrap(err, "unable to create up.sql")
	}

	if err := createFile(fmt.Sprintf("%s/down.sql", path)); err != nil {
		return errors.Wrap(err, "unable to create down.sql")
	}
	return nil
}

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
	return nil
}

func createFile(path string) error {
	if _, err := os.Create(path); err != nil {
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

			dirs = append(dirs, directory{
				name:      name,
				timestamp: 0,
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
