package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/app"
	"github.com/urfave/cli"
)

func makeApp(c *cli.Context) error {
	driver := c.String("driver")
	if driver == "" {
		return errors.New("driver is not defined")
	}

	a := app.New(driver, c.String("username"), c.String("password"), c.String("dbname"), c.String("host"), c.Int("port"))
	name := c.String("new")
	if name != "" {
		if err := a.Create(name); err != nil {
			return err
		}
	}

	run := c.String("run")
	if run != "" {
		if err := a.Run(run); err != nil {
			return err
		}
	}
	return nil
}
func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "driver",
			Value: "",
			Usage: "SQL driver",
		},
		cli.StringFlag{
			Name:  "new",
			Value: "",
			Usage: "Create a new migration",
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "password for db",
		},
		cli.StringFlag{
			Name:  "username",
			Value: "",
			Usage: "username for db",
		},
		cli.StringFlag{
			Name:  "dbname",
			Value: "",
			Usage: "db",
		},
		cli.StringFlag{
			Name:  "host",
			Value: "",
			Usage: "connect to db",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 0,
			Usage: "connect to db",
		},
	}
	app.Name = "migr"
	app.Usage = "tool for SQL migrations"
	app.Action = func(c *cli.Context) error {
		if err := makeApp(c); err != nil {
			log.Fatal(err)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
