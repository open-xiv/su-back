package user

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	rmongo "github.com/open-xiv/su-back/internal/repo/mongo"
	"github.com/open-xiv/su-back/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

type JwtCustomClaims struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
	jwt.RegisteredClaims
}

func Login(c echo.Context) error {
	// get name & password
	var person model.PersonInfo
	if err := c.Bind(&person); err != nil {
		zap.L().Debug("failed to bind person", zap.Error(err))
		return err
	}
	name := person.Name
	password := person.Password

	// check in mongo
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUserByName(coll, name)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if user.Person.Password != password {
		return echo.ErrUnauthorized
	}

	// set custom claims
	claims := &JwtCustomClaims{
		user.ID,
		user.Person.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token and send it as message
	t, err := token.SignedString([]byte(os.Getenv("SIGN_KEY")))
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, echo.Map{"token": t}, " ")
}

func LoginByKey(c echo.Context) error {
	// get name & key
	var person model.PersonInfo
	if err := c.Bind(&person); err != nil {
		zap.L().Debug("failed to bind person", zap.Error(err))
		return err
	}
	key := person.Key

	// check in mongo
	client := c.Get("mongo").(*mongo.Client)
	coll := client.Database("subook").Collection("users")
	user, err := rmongo.PullUserByKey(coll, key)
	if err != nil {
		return echo.ErrUnauthorized
	}

	// set custom claims
	claims := &JwtCustomClaims{
		user.ID,
		user.Person.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate encoded token and send it as message
	t, err := token.SignedString([]byte(os.Getenv("SIGN_KEY")))
	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, echo.Map{"token": t}, " ")
}

func CreateCustomClaims() echojwt.Config {
	c := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("SIGN_KEY")),
	}
	return c
}
