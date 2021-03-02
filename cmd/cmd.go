package cmd

import (
	"fmt"
	"strings"

	"github.com/Mayowa-Ojo/otter/internal"
	cli "github.com/urfave/cli/v2"
)

// Execute - main entry to cli
func Execute() *cli.App {
	app := &cli.App{
		Name:  "Escobar",
		Usage: "Take control of your heroku deployments",
		Commands: []*cli.Command{
			{
				Name:    "auth",
				Aliases: []string{"a"},
				Usage:   "authorize otter with your heroku account",
				Flags:   []cli.Flag{},
				Action: func(c *cli.Context) error {
					fmt.Println("Opening browser - authorize otter client with your heroku account. \nWaiting for authorization...")

					if err := AuthorizeClient(); err != nil {
						return cli.Exit(err.Error(), 1)
					}

					return nil
				},
			},
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Control your deployment's config vars",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "app",
						Aliases:  []string{"a"},
						Usage:    "your app name/id",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list all variables",
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "get variables from file",
					},
					&cli.StringFlag{
						Name:    "set",
						Aliases: []string{"s"},
						Usage:   "add a new variable(s)\nformat - key/value pairs delimited by ':'",
					},
					&cli.StringFlag{
						Name:    "remove",
						Aliases: []string{"r"},
						Usage:   "remove variable(s)",
					},
				},
				Action: func(c *cli.Context) error {
					app := c.String("app")
					spinner, err := internal.LoadingSpinner()
					spinner.Start()

					token, err := internal.GetAuthTokens()

					if err != nil {
						spinner.Prefix("something went wrong...")
						spinner.StopFail()
						return err
					}

					if c.Bool("list") {
						result, err := GetVariables(app, token)
						if err != nil {
							return err
						}

						fmt.Printf("data: %+v", result)
						spinner.Prefix("Done.")

						spinner.Stop()

						var out []interface{}

						for k, v := range result {
							out = append(out, map[string]interface{}{
								"key":   k,
								"value": v,
							})
						}

						table, err := internal.GenerateDataTable(out)
						if err != nil {
							return err
						}

						fmt.Println("we got here\n")
						fmt.Println(table.String())
						return nil
					}

					if file := c.String("file"); c.IsSet("file") {
						var source string

						if strings.Contains(file, "env") {
							source = "env"
						}
						if strings.Contains(file, "yaml") {
							source = "yaml"
						}
						if strings.Contains(file, "yml") {
							source = "yaml"
						}
						if strings.Contains(file, "json") {
							source = "json"
						}

						if err := UpsertVariables(app, token, file, source); err != nil {
							return err
						}

						spinner.Prefix("Done.")
						spinner.Stop()
						return nil
					}

					if variable := c.String("set"); c.IsSet("set") {
						key := strings.Split(variable, ":")[0]
						value := strings.Split(variable, ":")[1]
						kv := ConfigVar{
							key,
							value,
						}

						if err := UpsertVariable(app, token, kv); err != nil {
							return err
						}

						spinner.Prefix("Done.")
						spinner.Stop()
						return nil
					}

					if key := c.String("remove"); c.IsSet("remove") {
						if err := RemoveVariable(app, token, key); err != nil {
							return err
						}

						spinner.Prefix("Done.")
						spinner.Stop()
						return nil
					}

					return nil
				},
			},
		},
	}

	return app
}
