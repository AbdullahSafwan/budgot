package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Category struct {
	ent.Schema
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Enum("type").Values("income", "expense", "transfer").Immutable(),
		field.String("color").NotEmpty(),
		field.Bool("is_active").Default(true),
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

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Edges("owner").Unique(),
	}
}
