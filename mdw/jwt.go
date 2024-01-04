package mdw

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/exim-id/go-exim-base-lib/db"
	"github.com/exim-id/go-exim-base-lib/env"
	"github.com/exim-id/go-exim-base-lib/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Claims struct {
	jwt.StandardClaims
}

func GenerateToken() {

}

func GenerateEmptyToken() (string, time.Time, error) {
	exp := time.Now().Add(time.Minute * 30)
	claims := &Claims{StandardClaims: jwt.StandardClaims{ExpiresAt: exp.Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := token.SignedString([]byte(env.GetJwtSecret()))
	if err == nil {
		SaveToken(token)
	}
	return tokenString, exp, err
}

func SaveToken(token *jwt.Token) {
	db.NewMongoTransaction[bool](true).Transaction(func(d *mongo.Database) bool {
		claims := token.Claims.(*Claims)
		if err := claims.Valid(); err != nil {
			panic(err)
		}
		coll := d.Collection("sessions")
		session, exists := db.MongoDbFindFromCollection[models.Session](models.Session{Jti: claims.Id}, coll)
		if !exists {
			_, err := coll.InsertOne(context.Background(), models.Session{
				Jti:         claims.Id,
				Mode:        "anonym",
				CreatedDate: time.Now(),
			})
			if err != nil {
				panic(err)
			}
		} else {
			panic(fmt.Errorf("Session %s is exists", session.Id))
		}
		return false
	})
}

func setTokenCookies(name, token string, exp time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Expires = exp
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Name = name
	cookie.Value = token
	c.SetCookie(cookie)
}

func setUserCookies(employer *models.Employer, client *models.Client, exp time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Expires = exp
	cookie.Path = "/"
	if employer != nil {
		cookie.Name = "employer"
		cookie.Value = strings.Join([]string{employer.FrontName, employer.MidName, employer.LastName}, " ")
	} else if client != nil {
		cookie.Name = "client"
		cookie.Value = strings.Join([]string{client.FrontName, client.MidName, client.LastName, "@", client.Company.Name}, " ")
	}
	c.SetCookie(cookie)
}

func JWTErrorChecker(err error, c echo.Context) {
	c.Redirect(http.StatusUnauthorized, c.Echo().Reverse("error"))
}
