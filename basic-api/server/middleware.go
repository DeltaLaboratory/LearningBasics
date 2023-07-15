package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"basicapi/ent/user"
)

func (server *Server) authorized(c *fiber.Ctx) error {
	token, err := jwt.Parse(c.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		requestUser, err := server.db.User.Query().Where(user.UserID(claims["aud"].(string))).Only(c.Context())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
		c.Locals("user", requestUser)
		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid token")
	}
}
