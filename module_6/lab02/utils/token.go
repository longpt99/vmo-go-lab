package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func SignToken(id string) string {
	expirationTime := time.Now().Add(10 * 60 * time.Second)
	claims := &Claims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	return tokenString

}

func VerifyToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		if err.Error() == fmt.Sprintf("%s: %s", jwt.ErrTokenInvalidClaims.Error(), jwt.ErrTokenExpired.Error()) {
			return errors.New("token is expired")
		}
		return err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return errors.New("couldn't parse claims")
	}

	log.Printf("Token payload: %v\n", claims.Id)

	return nil

}
