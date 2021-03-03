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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "revoke",
						Aliases: []string{"r"},
						Usage:   "revoke your auth tokens",
					},
				},
				Action: func(c *cli.Context) error {
					spinner, err := internal.LoadingSpinner()
					if err != nil {
						spinner.Prefix("something went wrong...")
						spinner.StopFail()

						return cli.Exit(err.Error(), 1)
					}

					if c.IsSet("revoke") {
						spinner.Start()
						if err := internal.RevokeAuthorization(); err != nil {
							return cli.Exit(err.Error(), 1)
						}

						spinner.Prefix("Done")
						return nil
					}

					fmt.Println("Opening browser - authorize otter client with your heroku account.")
					spinner.Prefix("Waiting for authorization...")
					spinner.Start()
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

					tokens, err := internal.GetAuthTokens()

					if err != nil {
						spinner.Prefix("something went wrong...")
						spinner.StopFail()
						return err
					}

					if c.Bool("list") {
						result, err := GetVariables(app, tokens.AccessToken)
						if err != nil {
							return err
						}

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

						if err := UpsertVariables(app, tokens.AccessToken, file, source); err != nil {
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

						if err := UpsertVariable(app, tokens.AccessToken, kv); err != nil {
							return err
						}

						spinner.Prefix("Done.")
						spinner.Stop()
						return nil
					}

					if key := c.String("remove"); c.IsSet("remove") {
						if err := RemoveVariable(app, tokens.AccessToken, key); err != nil {
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
