package utils

import (
	"fmt"
	"github.com/HoseonYim/isfree-backend/configs"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthTokenClaims struct {
	UserID string `json:"id"` // 유저 ID
	jwt.RegisteredClaims
}

func CreateToken(userID string) (string, error) {
	var err error
	//Creating Access Token
	now := time.Now()
	at := AuthTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Duration(config.JWT_EXPIRE) * time.Second)},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, at)
	tokenString, err := token.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ParseToken(tokenString string) (*AuthTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET), nil
	})
	if err != nil {
		fmt.Println(err)
		return &AuthTokenClaims{}, err
	}
	claims, ok := token.Claims.(*AuthTokenClaims)
	if !ok {
		return &AuthTokenClaims{}, err
	}
	return claims, nil
}

var config = configs.Config
