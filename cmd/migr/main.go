package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"log"
	"os"
)

func makeApp(c *cli.Context) error {
	driver := c.String("driver")
	if driver == "" {
		return errors.New("driver is not defined")
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
