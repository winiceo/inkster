package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"

	serverApp "github.com/dominik-zeglen/inkster/app"
	"github.com/dominik-zeglen/inkster/core"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			operation := c.Args().Get(0)
			if operation == "runserver" {
				serverApp.Run()
				return nil
			}

			if operation == "add-user" {
				email := c.Args().Get(1)
				password := c.Args().Get(2)
				newUser := core.User{
					Email:  email,
					Active: true,
				}
				err := newUser.CreatePassword(password)
				if err != nil {
					return err
				}

				_, err = serverApp.DataSource.AddUser(newUser)
				if err != nil {
					return err
				}
				fmt.Println("Added user " + email)
				return nil
			}
		}
		fmt.Println("No operation given")
		return nil

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
