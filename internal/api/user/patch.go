package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/internal/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// check token
	uToken := c.Get("user").(*jwt.Token)
	claims := uToken.Claims.(*JwtCustomClaims)
	uId := claims.ID
	if uId != id {
		zap.L().Debug("permission denied (token != id)")
		return echo.ErrUnauthorized
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
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("tale").Collection("users")
	err = rmongo.PatchUser(coll, id, m)
	if err != nil {
		zap.L().Debug("failed to patch user", zap.Error(err))
		return err
	}

	// return {"success": bool}
	return c.JSONPretty(http.StatusOK, map[string]bool{"success": true}, "  ")
}

func _(c echo.Context) error {
	// get user id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		zap.L().Debug("failed to parse id", zap.Error(err))
		return err
	}

	// get fight id
	fId, err := primitive.ObjectIDFromHex(c.Param("fight_id"))
	if err != nil {
		zap.L().Debug("failed to parse fight id", zap.Error(err))
		return err
	}

	// check token
	uToken := c.Get("user").(*jwt.Token)
	claims := uToken.Claims.(*JwtCustomClaims)
	uId := claims.ID
	if uId != id {
		zap.L().Debug("permission denied (token != id)")
		return echo.ErrUnauthorized
	}

	// mongo
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("tale").Collection("users")
	err = rmongo.InsertFight(coll, id, fId)
	if err != nil {
		zap.L().Debug("failed to insert fight", zap.Error(err))
		return err
	}

	// return {"success": bool}
	return c.JSONPretty(http.StatusOK, map[string]bool{"success": true}, "  ")
}
