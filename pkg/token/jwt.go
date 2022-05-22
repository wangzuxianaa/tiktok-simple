package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("douyin")

type Claims struct {
	UserId   int64
	Username string
	jwt.StandardClaims
}

func GenerateToken(id int64, name string) (string, error) {
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

func ParseToken(tokenStr string) (*Claims, bool) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, true
	}
	return nil, false
}
