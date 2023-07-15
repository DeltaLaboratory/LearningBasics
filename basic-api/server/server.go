package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/goccy/go-json"

	"basicapi/ent"
)

const JWT_SECRET = "pepsi zero sugar"

type Server struct {
	app *fiber.App

	db *ent.Client
}

func NewServer(db *ent.Client) *Server {
	return &Server{
		app: fiber.New(fiber.Config{
			JSONDecoder: json.Unmarshal,
			JSONEncoder: json.Marshal,
		}),
		db: db,
	}
}

func (server *Server) Run(addr string) error {
	server.route()
	return server.app.Listen(addr)
}
