package fight

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/pkg/model"
	"go.uber.org/zap"
	"net/http"
)

func Push(c echo.Context) error {
	// bind fight
	var fight model.Fight
	if err := c.Bind(&fight); err != nil {
		zap.L().Debug("failed to bind fight", zap.Error(err))
		return err
	}
	fight.ServerRecord.IP = c.RealIP()

	// mongo
	client := config.MongoClient
	coll := client.Database("subook").Collection("fights")
	family, err := rmongo.PushFight(coll, fight)
	if err != nil {
		zap.L().Debug("failed to push fight", zap.Error(err))
		return err
	}

	// return fight
	return c.JSONPretty(http.StatusOK, family, "  ")
}
