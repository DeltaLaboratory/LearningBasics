package main

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"basicapi/ent"
	"basicapi/ent/article"
	"basicapi/ent/comment"
	"basicapi/ent/user"
	tool "basicapi/internal"
)

const JWT_SECRET = "pepsi zero sugar"

type Server struct {
	app *fiber.App

	db *ent.Client
}

func Authorized(c *fiber.Ctx) error {
	server := c.Locals("server").(*Server)

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

func main() {
	var err error
	server := Server{
		app: fiber.New(),
	}
	server.db, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		panic(err)
	}

	_ = server.db.Schema.Create(context.Background())

	server.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	server.app.Use(func(c *fiber.Ctx) error {
		c.Locals("server", &server)
		return c.Next()
	})

	server.app.Post("/account/register", func(ctx *fiber.Ctx) error {
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
	})

	server.app.Post("/account/login", func(ctx *fiber.Ctx) error {
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
	})

	server.app.Get("/articles", func(ctx *fiber.Ctx) error {
		var err error
		off := ctx.QueryInt("offset", 0)
		if off == -1 || err != nil {
			off = 0
		}

		fmt.Println(off)

		articles, err := server.db.Article.Query().Order(article.ByID(sql.OrderDesc())).Offset(off).Limit(20).All(ctx.Context())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.Status(fiber.StatusOK).JSON(ArticleListResponse{
			Articles: articles,
		})
	})

	server.app.Post("/articles", Authorized, func(ctx *fiber.Ctx) error {
		request := new(NewArticleRequest)
		if err := ctx.BodyParser(&request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		user := (ctx.Locals("user")).(*ent.User)
		if article, err := server.db.Article.Create().SetTitle(request.Title).SetContent(request.Content).SetAuthor(user).Save(ctx.Context()); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		} else {
			return ctx.Status(fiber.StatusOK).JSON(NewArticleResponse{
				Article: article,
			})
		}
	})

	server.app.Get("/articles/:id", func(ctx *fiber.Ctx) error {
		articleId, err := ctx.ParamsInt("id")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		article, err := server.db.Article.Get(ctx.Context(), articleId)
		if err != nil {
			if _, ok := err.(*ent.NotFoundError); ok {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			} else {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}

		return ctx.Status(fiber.StatusOK).JSON(GetArticleResponse{
			Article: article,
		})
	})

	server.app.Get("/articles/:id/comments", func(ctx *fiber.Ctx) error {
		var err error
		off := ctx.QueryInt("offset", 0)

		var articleId int
		if articleId, err = ctx.ParamsInt("id"); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var comments []*ent.Comment
		if comments, err = server.db.Comment.Query().
			WithArticle(func(aq *ent.ArticleQuery) {
				aq.Where(article.ID(articleId))
			}).
			Order(comment.ByID(sql.OrderDesc())).
			Offset(off).
			Limit(20).
			All(ctx.Context()); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.Status(fiber.StatusOK).JSON(CommentListResponse{
			Comments: comments,
		})
	})

	server.app.Post("/articles/:id/comments", Authorized, func(ctx *fiber.Ctx) error {
		var articleId int
		if articleId, err = ctx.ParamsInt("id"); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		request := new(NewCommentRequest)
		if err := ctx.BodyParser(&request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		user := ctx.Locals("user").(*ent.User)
		if comment, err := server.db.Comment.Create().SetAuthor(user).SetContent(request.Content).SetArticleID(articleId).Save(ctx.Context()); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		} else {
			return ctx.Status(fiber.StatusOK).JSON(NewCommentResponse{
				Comment: comment,
			})
		}

	})

	server.app.Listen(":80")
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type NewArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NewArticleResponse struct {
	Article *ent.Article `json:"article"`
}

type ArticleListResponse struct {
	Articles []*ent.Article `json:"articles"`
}

type GetArticleResponse struct {
	Article *ent.Article `json:"article"`
}

type CommentListResponse struct {
	Comments []*ent.Comment `json:"comments"`
}

type NewCommentRequest struct {
	Content string `json:"content"`
}

type NewCommentResponse struct {
	Comment *ent.Comment `json:"comment"`
}
