package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"basicapi/ent/user"
	tool "basicapi/internal"
)

func (server *Server) register(ctx *fiber.Ctx) error {
	request := new(RegisterRequest)
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	hashedPassword, err := tool.Generate([]byte(request.Password))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if user, err := server.db.User.Create().
		SetUserID(request.Username).
		SetPassword(hashedPassword).
		Save(ctx.Context()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	} else {
		return ctx.Status(fiber.StatusOK).JSON(user)
	}
}

func (server *Server) login(ctx *fiber.Ctx) error {
	request := new(RegisterRequest)
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if user, err := server.db.User.Query().
		Where(user.UserID(request.Username)).
		Only(ctx.Context()); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
	} else {
		if ok, err := tool.Verify([]byte(request.Password), user.Password); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		} else {
			if ok {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"nbf": time.Now().Unix(),
					"exp": time.Now().Add(time.Hour * 24).Unix(),
					"aud": user.UserID,
				})
				tokenString, err := token.SignedString([]byte(JWT_SECRET))
				if err != nil {
					return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}
				return ctx.Status(fiber.StatusOK).JSON(LoginResponse{
					Token: tokenString,
				})
			} else {
				return ctx.Status(fiber.StatusUnauthorized).SendString("wrong password")
			}
		}
	}

}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
