package httpserver

import (
	"context"
	"crou-api/config"
	"crou-api/framework/httpserver/router"
	app "crou-api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"go.uber.org/fx"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

const LOG_FORMAT = "[${time}] ${ip} - ${status} ${method} ${path} ${queryParams} ${body} - ${latency} \n"

func Api(
	lc fx.Lifecycle,
	conf *config.Config,
	stx *app.ServiceContext,
) *fiber.App {

	app := fiber.New(fiber.Config{
		ErrorHandler: defaultErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(conf.Cors.Origins, ","),
		AllowMethods:     strings.Join(conf.Cors.Methods, ","),
		AllowHeaders:     strings.Join(conf.Cors.Headers, ","),
		AllowCredentials: conf.Cors.Credentials,
	}))

	app.Use(requestid.New(requestid.Config{
		Next:       nil,
		Header:     fiber.HeaderXRequestID,
		Generator:  utils.UUID,
		ContextKey: "requestid",
	}))
	app.Use(logger.New(logger.Config{
		Format:     LOG_FORMAT,
		TimeFormat: "2006-01-02 15:04:05",
	}))
	if conf.Log.Type == "file" {
		app.Use(logger.New(logger.Config{
			Format:        LOG_FORMAT,
			TimeFormat:    "2006-01-02 15:04:05",
			DisableColors: true,
			Output: &lumberjack.Logger{
				Filename:   conf.Log.FileName,
				MaxSize:    20, // megabytes
				MaxBackups: 3,
				MaxAge:     28, // days
			},
		}))
	}
	app.Use(recover.New())

	router.Route(conf, app, stx)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(conf.Host + ":" + conf.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
	return app
}
