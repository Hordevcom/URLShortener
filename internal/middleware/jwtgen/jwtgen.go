package jwtgen

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// type JWTGen struct {
// 	token_exp  time.Duration
// 	secret_key string
// }

// func NewJWTGen() JWTGen {
// 	return JWTGen{
// 		token_exp:  time.Hour * 12,
// 		secret_key: "supersecretkey",
// 	}
// }

const TOKEN_EXP = time.Hour * 12
const SECRET_KEY = "supersecretkey"

func BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: 2,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}

	return claims.UserID
}
