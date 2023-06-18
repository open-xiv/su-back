package user

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/pkg/model"
	"go.uber.org/zap"
	"net/http"
)

func Init(c echo.Context) error {
	// bind user
	var user model.User
	if err := c.Bind(&user); err != nil {
		zap.L().Debug("failed to bind user", zap.Error(err))
		return err
	}
	user.ServerRecord.IP = c.RealIP()

	// mongo
	client := config.MongoClient
	coll := client.Database("subook").Collection("users")
	id, err := rmongo.InitUser(coll, user)
	if err != nil {
		zap.L().Debug("failed to init user", zap.Error(err))
		return err
	}

	// return {"id": id}
	return c.JSONPretty(http.StatusOK, map[string]interface{}{"id": id}, "  ")
}
