package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Category struct {
	ent.Schema
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Enum("type").Values("income", "expense", "transfer"),
		field.String("color").NotEmpty(),
	}
}

func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("categories").Unique().Required(),
		edge.To("transactions", Transaction.Type).
			StorageKey(edge.Column("category_id")),
		edge.To("budgets", Budget.Type).
			StorageKey(edge.Column("category_id")),
	}
}
