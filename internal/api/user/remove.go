package user

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
)

func Remove(c echo.Context) error {
	// get user id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		zap.L().Debug("failed to parse id", zap.Error(err))
		return err
	}

	// mongo
	client := config.MongoClient
	coll := client.Database("tale").Collection("users")
	if err := rmongo.RemoveUser(coll, id); err != nil {
		zap.L().Debug("failed to remove user", zap.Error(err))
		return err
	}

	// return {"success": bool}
	return c.JSONPretty(http.StatusOK, map[string]bool{"success": true}, "  ")
}
