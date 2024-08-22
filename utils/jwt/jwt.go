package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AppConfig struct {
	JWT_SECRET_KEY string
}

func loadEnv() *AppConfig {
	var res = new(AppConfig)
	godotenv.Load(".env")

	// if err != nil {
	// 	logrus.Error("Config : Cannot load config file,", err.Error())
	// 	return nil
	// }

	if secretKey, found := os.LookupEnv("JWT_SECRET_KEY"); found {
		res.JWT_SECRET_KEY = secretKey
	}
	return res
}

func JWTMiddleware() echo.MiddlewareFunc {
	key := loadEnv()
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(key.JWT_SECRET_KEY),
		SigningMethod: "HS256",
	})
}

func CreateToken(userId string, role string, religion string) (string,error) {
	key := loadEnv()
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["role"] = role
	claims["religion"] = religion
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key.JWT_SECRET_KEY))

}

func ExtractTokenUserId(c echo.Context) (string, string, string, error) {

	user := c.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		Id := claims["userId"].(string)
		Role := claims["role"].(string)
		Religion := claims["religion"].(string)
		return Id, Role, Religion, nil
	}
	return "", "", "", errors.New("invalid token")
}
