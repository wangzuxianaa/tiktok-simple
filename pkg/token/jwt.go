package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("douyin")

type Claims struct {
	UserId   uint
	Username string
	jwt.StandardClaims
}

func GenerateToken(id uint, name string) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := Claims{
		UserId:   id,
		Username: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tiktok",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)

	return tokenStr, err
}
