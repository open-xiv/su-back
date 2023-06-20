package fight

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/internal/tools"
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

	// check token
	uToken := c.Get("user").(*jwt.Token)
	claims := uToken.Claims.(*tools.JwtCustomClaims)
	uId := claims.ID
	if uId != fight.UserID {
		zap.L().Debug("permission denied (token != id)")
		return echo.ErrUnauthorized
	}

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
