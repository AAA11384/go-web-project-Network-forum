package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpireDuration = time.Hour * 2

var key = []byte("Jialin Liang")

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"password"`
	jwt.RegisteredClaims
}

func GenToken(userId int64, username string) (string, error) {
	c := MyClaims{
		userId,
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my_project",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(key)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var rc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, rc, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return rc, nil
	}
	return nil, err
}
