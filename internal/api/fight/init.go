package fight

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/pkg/model"
	"go.uber.org/zap"
	"net/http"
)

func Init(c echo.Context) error {
	// fight record json -> bind fight
	var fight model.Fight
	if err := c.Bind(&fight.FightRecord); err != nil {
		zap.L().Debug("failed to bind fight", zap.Error(err))
		return err
	}

	// mongo
	client := config.MongoClient
	coll := client.Database("subook").Collection("fights")
	id, err := rmongo.InitFights(coll, fight)
	if err != nil {
		zap.L().Debug("failed to init fight", zap.Error(err))
		return err
	}

	// return {"id": id}
	return c.JSONPretty(http.StatusOK, map[string]interface{}{"id": id}, "  ")
}
