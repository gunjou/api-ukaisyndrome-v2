package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(rdb *redis.Client, jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(ErrorResponse{Message: "missing token"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(ErrorResponse{Message: "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(ErrorResponse{
				Message: "invalid claims",
			})
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			return c.Status(401).JSON(ErrorResponse{
				Message: "invalid token (sub missing)",
			})
		}

		userID, err := strconv.Atoi(sub)
		if err != nil {
			return c.Status(401).JSON(ErrorResponse{
				Message: "invalid user id",
			})
		}

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(401).JSON(ErrorResponse{
				Message: "invalid role",
			})
		}

		platform, _ := claims["platform"].(string)

		// 🔥 FIX KEY
		key := "session:" + role + ":" + platform + ":" + strconv.Itoa(userID)

		exists, err := rdb.Exists(c.Context(), key).Result()
		if err != nil || exists == 0 {
			return c.Status(401).JSON(ErrorResponse{Message: "session expired"})
		}

		c.Locals("sub", userID)
		c.Locals("role", role)
		c.Locals("platform", platform)

		return c.Next()
	}
}