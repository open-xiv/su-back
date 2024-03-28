package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
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
	//e.Debug = true
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
	// cors
	e.Use(middleware.CORS())
	// recover
	e.Use(middleware.Recover())

	// connect to mongo
	mongoClient := config.ConnectDB()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mongo", mongoClient)
			return next(c)
		}
	})

	// restful api (golang backend)
	// api (backend)
	b := e.Group("")

	// public api
	pub := b.Group("/public")
	// server status
	pub.GET("/status", server.Status)
	// user fights record
	pub.GET("/user/:name/fights", user.PullRecords)
	pub.GET("/user/:name/meta", user.PullMeta)
	pub.GET("/user/:name/avatar", user.PullAvatar)
	// fight record
	pub.GET("/fight/:id", fight.Pull)

	// protect api
	pro := b.Group("/protect")
	pro.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	// user
	pro.POST("/user", user.Init)   // create user
	pro.POST("/login", user.Login) // login -> token

	// private api
	pri := b.Group("/private")
	pri.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	pri.Use(echojwt.WithConfig(user.CreateCustomClaims()))
	// user
	priUser := pri.Group("/user")
	priUser.GET("", user.Pull)
	priUser.GET("/:id", user.PullByID)
	priUser.PUT("/:id", user.Push)
	priUser.DELETE("/:id", user.Remove)
	priUser.PATCH("/:id", user.Patch)
	// fight
	priFight := pri.Group("/fight")
	priFight.POST("", fight.Init)
	priFight.PUT("/:id", fight.Push)
	priFight.DELETE("/:id", fight.Remove)

	// echo server
	e.Logger.Fatal(e.Start(":8123"))
}
