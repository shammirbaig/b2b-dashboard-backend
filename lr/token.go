package lr

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/savsgio/atreugo/v11"
	"log"
	"os"
	"time"
)

func CreateSessionToken(userId string) string {
	// Define the token claims
	claims := jwt.MapClaims{
		"uid": userId,
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}

	return tokenString
}

func ValidateToken(ctx *atreugo.RequestCtx, tokenString string) error {
	if tokenString == "" {
		return ctx.ErrorResponse(fmt.Errorf("missing token"), 401)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return ctx.ErrorResponse(fmt.Errorf("invalid token"), 401)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.ErrorResponse(fmt.Errorf("invalid claims"), 401)
	}

	ctx.SetUserValue("uid", claims["uid"])

	return nil
}
