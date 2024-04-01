package user

import (
	"crypto/md5"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
)

func Pull(c echo.Context) error {
	// check token
	uToken := c.Get("user").(*jwt.Token)
	claims := uToken.Claims.(*JwtCustomClaims)
	uId := claims.ID

	// mongo
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUser(coll, uId)
	if err != nil {
		zap.L().Debug("failed to pull user", zap.Error(err))
		return err
	}

	// return user
	return c.JSONPretty(http.StatusOK, user, "  ")
}

func PullByID(c echo.Context) error {
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

	// mongo
	client := c.Get("mongo").(*mongo.Client)
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
	// get username
	name := c.Param("name")

	// mongo
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUserByName(coll, name)
	if err != nil {
		zap.L().Debug("failed to pull user", zap.Error(err))
		return err
	}

	// return user fights
	return c.JSONPretty(http.StatusOK, echo.Map{"fight_ids": user.FightIDs}, "  ")
}

func PullMeta(c echo.Context) error {
	// get username
	name := c.Param("name")

	// mongo
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUserByName(coll, name)
	if err != nil {
		zap.L().Debug("failed to pull user", zap.Error(err))
		return err
	}

	// return user meta
	return c.JSONPretty(http.StatusOK, echo.Map{"meta": user.Meta}, "  ")
}

func PullAvatar(c echo.Context) error {
	// get username
	name := c.Param("name")

	// mongo
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUserByName(coll, name)
	if err != nil {
		zap.L().Debug("failed to pull user", zap.Error(err))
		return err
	}

	// if no avatar
	avatarURL := user.Person.AvatarURL
	if avatarURL == "" {
		// use gravatar (email md5)
		avatarURL = "https://www.gravatar.com/avatar/" + fmt.Sprintf("%x", md5.Sum([]byte(user.Person.Email)))
	}

	// return user meta
	return c.JSONPretty(http.StatusOK, echo.Map{"avatar_url": avatarURL}, "  ")
}
