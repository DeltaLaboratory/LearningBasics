package server

import "github.com/gofiber/fiber/v2/middleware/cors"

func (server *Server) route() {
	server.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	account := server.app.Group("/account")
	{
		account.Post("/register", server.register)
		account.Post("/login", server.login)
	}

	articles := server.app.Group("/articles")
	{
		articles.Get("/", server.listArticle)
		articles.Get("/:id", server.getArticle)
		articles.Get("/:id/comments", server.listComment)

		authorized := articles.Group("")
		authorized.Use(server.authorized)
		{
			authorized.Post("/", server.createArticle)
			authorized.Post("/:id/comments", server.createComment)
		}
	}
}
