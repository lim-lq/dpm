package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var dpmSecret = []byte("__DPM_IS_UNDER_DEVELOPING__")

func NewToken(user string) (string, error) {
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"user":        user,
	// 	"expiredTime": time.Now().Unix(),
	// })
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 30)),
		Issuer:    "DPM",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(dpmSecret)
}

func CheckToken(tokenString string) error {
	tokenStr := strings.TrimSpace(tokenString)
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return dpmSecret, nil
	})
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		fmt.Println(claims.Issuer)
	} else {
		return err
	}

	return nil
}
