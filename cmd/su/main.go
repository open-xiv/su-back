package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/open-xiv/su-back/config"
	"github.com/open-xiv/su-back/internal/api/fight"
	"github.com/open-xiv/su-back/internal/api/server"
	"github.com/open-xiv/su-back/internal/api/user"
	"go.uber.org/zap"
)

func main() {
	// echo root
	e := echo.New()
	e.Debug = true
	e.HideBanner = true

	// logger
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			zap.L().Debug("failed to sync logger", zap.Error(err))
		}
	}(logger)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))
	zap.ReplaceGlobals(logger)

	// middleware
	// remove trailing slash
	e.Pre(middleware.RemoveTrailingSlash())
	// rate limit
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	// cors
	e.Use(middleware.CORS())

	// connect to mongo
	config.ConnectDB()

	// static files (vue frontend)
	e.Use(middleware.Static("./web"))

	// restful api (golang backend)
	b := e.Group("/api")
	// server status
	b.GET("/status", server.Status)
	// user
	e.POST("/user", user.Init)
	e.GET("/user/:id", user.Pull)
	e.PUT("/user/:id", user.Push)
	e.DELETE("/user/:id", user.Remove)
	e.PATCH("/user/:id", user.Patch)
	// fight
	e.POST("/fight", fight.Init)
	e.GET("/fight/:id", fight.Pull)
	e.PUT("/fight/:id", fight.Push)
	e.DELETE("/fight/:id", fight.Remove)

	// echo server
	e.Logger.Fatal(e.Start(":8123"))
}
