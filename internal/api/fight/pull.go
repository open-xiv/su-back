package fight

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
)

func Pull(c echo.Context) error {
	// get fight id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		zap.L().Debug("failed to parse id", zap.Error(err))
		return err
	}

	// mongo
	client := config.MongoClient
	coll := client.Database("subook").Collection("fights")
	fight, err := rmongo.PullFight(coll, id)
	if err != nil {
		zap.L().Debug("failed to pull fight", zap.Error(err))
		return err
	}

	// return fight
	return c.JSONPretty(http.StatusOK, fight, "  ")
}
