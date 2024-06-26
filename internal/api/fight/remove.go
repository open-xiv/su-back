package fight

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/internal/api/user"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
)

func Remove(c echo.Context) error {
	// get fight id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		zap.L().Debug("failed to parse id", zap.Error(err))
		return err
	}

	// query fight
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("subook").Collection("fights")
	fight, err := rmongo.PullFight(coll, id)
	if err != nil {
		zap.L().Debug("failed to query fight", zap.Error(err))
		return err
	}

	// check token
	uToken := c.Get("user").(*jwt.Token)
	claims := uToken.Claims.(*user.JwtCustomClaims)
	uId := claims.ID
	if uId != fight.UserID {
		zap.L().Debug("permission denied (token != id)")
		return echo.ErrUnauthorized
	}

	// mongo
	if err := rmongo.RemoveFight(coll, id); err != nil {
		zap.L().Debug("failed to remove fight", zap.Error(err))
		return err
	}

	// return {"success": bool}
	return c.JSONPretty(http.StatusOK, map[string]bool{"success": true}, "  ")
}
