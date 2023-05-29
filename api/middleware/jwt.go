package middleware

import (
	"fmt"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]

		if !ok {
			fmt.Println("token not present in the header")
			return api.ErrUnauthorized()

		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		// check if token is okai
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		// Check if the current time is after the expiration time
		currentTime := time.Now().UTC().Unix()
		if currentTime > expires {
			return api.NewError(http.StatusUnauthorized, "token expired")
		}

		userID := claims["id"].(string)

		user, err := userStore.GetUserByID(c.Context(), userID)

		if err != nil {
			return api.ErrUnauthorized()
		}
		// Set the current authenticated user to the context value
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, api.ErrUnauthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, api.ErrUnauthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, api.ErrUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, api.ErrUnauthorized()
	}

	return claims, nil
}
