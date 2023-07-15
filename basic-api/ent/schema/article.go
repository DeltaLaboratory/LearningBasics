package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("content").NotEmpty().Sensitive(),

		field.Time("created_at").Default(time.Now),

		field.Int("author_id"),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).
			Field("author_id").
			Ref("articles").
			Required().
			Unique(),
		edge.To("comments", Comment.Type),
	}
}
