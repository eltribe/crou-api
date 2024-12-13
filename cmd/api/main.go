package main

import (
	"crou-api/config"
	"crou-api/config/database"
	"crou-api/framework/httpserver"
	app "crou-api/internal"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"log"
	"os"
)

// @title     CROU API documentation
// @version   1.0.0
// @BasePath  /
func main() {
	app := Api()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Api() *cli.App {
	app := &cli.App{
		Name:    Name,
		Usage:   Usage,
		Version: Version,
		Flags:   Flags,
		Action: func(c *cli.Context) error {
			// Config Load
			filename := c.String(EnvFlag.Name)
			conf := config.LoadConfigFile(filename)
			conf.Log.FileName = "pliper-api.log"

			fx.New(
				fx.Provide(
					func() *config.Config { return conf },
					database.NewDatabase,
					config.NewOauth,
				),
				app.Ctx,
				fx.Invoke(
					database.AutoMigration,
					httpserver.Api,
				),
			).Run()
			return nil
		},
	}
	return app
}
