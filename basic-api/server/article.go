package server

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2"

	"basicapi/ent"
	"basicapi/ent/article"
	"basicapi/ent/comment"
	tool "basicapi/internal"
)

func (server *Server) listArticle(ctx *fiber.Ctx) error {
	off := ctx.QueryInt("offset", 0)
	if off < 0 {
		off = 0
	}

	articles, err := server.db.Article.Query().
		Order(article.ByID(sql.OrderDesc())).
		Offset(off).
		Limit(20).
		All(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	articleList := new(ArticleListResponse)
	articleList.Articles = []ArticleEntry{}

	for _, atcl := range articles {
		author, err := atcl.QueryAuthor().Only(ctx.Context())
		if err != nil {
			articleList.Articles = append(articleList.Articles, ArticleEntry{
				ID:      atcl.ID,
				Title:   atcl.Title,
				Content: tool.Cut(atcl.Content, 25),
				Created: atcl.CreatedAt,
				Author:  "Ghost",
			})
			continue
		}
		articleList.Articles = append(articleList.Articles, ArticleEntry{
			ID:      atcl.ID,
			Title:   atcl.Title,
			Content: tool.Cut(atcl.Content, 25),
			Created: atcl.CreatedAt,
			Author:  author.UserID,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(articleList)
}

func (server *Server) createArticle(ctx *fiber.Ctx) error {
	request := new(NewArticleRequest)
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user := (ctx.Locals("user")).(*ent.User)
	if article, err := server.db.Article.Create().
		SetTitle(request.Title).
		SetContent(request.Content).
		SetAuthor(user).
		Save(ctx.Context()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	} else {
		return ctx.Status(fiber.StatusOK).JSON(NewArticleResponse{
			Article: article,
		})
	}
}

func (server *Server) getArticle(ctx *fiber.Ctx) error {
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

	author, err := article.QueryAuthor().Only(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(ArticleEntry{
		ID:      article.ID,
		Title:   article.Title,
		Content: article.Content,
		Created: article.CreatedAt,
		Author:  author.UserID,
	})
}

func (server *Server) listComment(ctx *fiber.Ctx) error {
	var err error
	off := ctx.QueryInt("offset", 0)
	if off < 0 {
		off = 0
	}

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
}

func (server *Server) createComment(ctx *fiber.Ctx) error {
	var err error

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
}

type NewArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleEntry struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
}

type NewArticleResponse struct {
	Article *ent.Article `json:"article"`
}

type ArticleListResponse struct {
	Articles []ArticleEntry `json:"articles"`
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
