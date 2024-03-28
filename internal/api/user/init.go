package user

import (
	"github.com/labstack/echo/v4"
	"github.com/open-xiv/su-back/config"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
)

func Init(c echo.Context) error {
	// bind user
	var user model.User
	if err := c.Bind(&user.Person); err != nil {
		zap.L().Debug("failed to bind user", zap.Error(err))
		return err
	}
	user.ServerRecord.IP = c.RealIP()

	// check empty key
	for _, x := range []string{user.Person.Name, user.Person.Email, user.Person.Password} {
		if x == "" {
			zap.L().Debug("empty key")
			return echo.ErrBadRequest
		}
	}

	// reset fights id & meta
	user.FightIDs = []primitive.ObjectID{}
	user.Meta.Total = 2000

	// mongo
	client := config.MongoClient
	coll := client.Database("subook").Collection("users")
	id, err := rmongo.InitUser(coll, user)
	if err != nil {
		zap.L().Debug("failed to init user", zap.Error(err))
		if err.Error() == "user name already exists" {
			return c.JSONPretty(http.StatusForbidden, echo.Map{"error": err.Error()}, "  ")
		}
		return err
	}

	// return {"id": id}
	return c.JSONPretty(http.StatusOK, echo.Map{"id": id}, "  ")
}
