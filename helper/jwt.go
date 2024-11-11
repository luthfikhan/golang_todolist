package helper

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func CreateToken() (string, error) {
	claims := &Claims{
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(authHeader string) (*Claims, error) {
	re := regexp.MustCompile(`^Bearer (\S+)$`)
	matches := re.FindStringSubmatch(authHeader)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid Authorization header format")
	}

	tokenString := matches[1]

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
