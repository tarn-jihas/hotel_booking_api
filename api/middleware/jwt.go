package middleware

import (
	"fmt"
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
			return fmt.Errorf("unauthorized")
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
			return fmt.Errorf("token has expired")
		}

		userID := claims["id"].(string)

		user, err := userStore.GetUserByID(c.Context(), userID)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(genericResp{
				Type: "error",
				Msg:  "unauthorized",
			})
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
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
