package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/internal/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
)

func Pull(c echo.Context) error {
	// get user id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		zap.L().Debug("failed to parse id", zap.Error(err))
		return err
	}

	// check token
	uToken := c.Get("user").(*jwt.Token)
	claims := uToken.Claims.(*tools.JwtCustomClaims)
	uId := claims.ID
	if uId != id {
		zap.L().Debug("permission denied (token != id)")
		return echo.ErrUnauthorized
	}

	// mongo
	client := config.MongoClient
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUser(coll, id)
	if err != nil {
		zap.L().Debug("failed to pull user", zap.Error(err))
		return err
	}

	// return user
	return c.JSONPretty(http.StatusOK, user, "  ")
}

func PullRecords(c echo.Context) error {
	// get user name
	name := c.Param("name")

	// mongo
	client := config.MongoClient
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUserByName(coll, name)
	if err != nil {
		zap.L().Debug("failed to pull user", zap.Error(err))
		return err
	}

	// return user
	return c.JSONPretty(http.StatusOK, echo.Map{"fight_ids": user.FightIDs}, "  ")
}
