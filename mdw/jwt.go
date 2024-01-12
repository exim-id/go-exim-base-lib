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

func RefreshToken(tokenStr string) (string, error) {
	claims, err := ParseToken(tokenStr)
	if err != nil {
		return "", err
	}
	return GenerateToken(claims)
}

func ParseToken(tokenStr string) (Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Invalid token")
		} else if method != jwt.SigningMethodRS512 {
			return nil, fmt.Errorf("Method not equals")
		}
		return []byte(env.GetJwtSecret()), nil
	})
	if err != nil {
		return Claims{}, err
	}
	claims, ok := token.Claims.(Claims)
	if !ok || !token.Valid {
		return Claims{}, fmt.Errorf("Invalid token")
	}
	return claims, err
}

func GenerateEmptyToken() (string, error) {
	exp := time.Now().Add(time.Minute * 30)
	claims := &Claims{StandardClaims: jwt.StandardClaims{ExpiresAt: exp.Unix()}}
	return GenerateToken(*claims)
}

func GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := token.SignedString([]byte(env.GetJwtSecret()))
	if err == nil {
		SaveToken(token)
	}
	return tokenString, err
}

func SaveToken(token *jwt.Token) error {
	return db.NewMongoTransaction[error](true).Transaction(func(d *mongo.Database) error {
		claims := token.Claims.(*Claims)
		if err := claims.Valid(); err != nil {
			return err
		}
		coll := d.Collection("sessions")
		session, exists := db.MongoDbFindFromCollection[models.Session](models.Session{Jti: claims.Id}, coll)
		if !exists {
			ctx, _ := context.WithTimeout(context.Background(), time.Minute*30)
			_, err := coll.InsertOne(ctx, models.Session{
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
		return nil
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
