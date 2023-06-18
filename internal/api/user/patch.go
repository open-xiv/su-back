package user

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/internal/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
)

func Patch(c echo.Context) error {
	// get user id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		zap.L().Debug("failed to parse id", zap.Error(err))
		return err
	}

	// bind user
	m := make(map[string]interface{})
	if err := c.Bind(&m); err != nil {
		zap.L().Debug("failed to bind user", zap.Error(err))
		return err
	}
	m = tools.Flatten(m)
	m["record.ip"] = c.RealIP()
	zap.L().Debug("patch user", zap.Any("user", m))

	// mongo
	client := config.MongoClient
	coll := client.Database("tale").Collection("users")
	err = rmongo.PatchUser(coll, id, m)
	if err != nil {
		zap.L().Debug("failed to patch user", zap.Error(err))
		return err
	}

	// return {"success": bool}
	return c.JSONPretty(http.StatusOK, map[string]bool{"success": true}, "  ")
}