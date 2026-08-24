package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Budget struct {
	ent.Schema
}

func (Budget) Fields() []ent.Field {
	return []ent.Field{
		field.Int("month").Min(1).Max(12),
		field.Int("year"),
		field.Int64("amount").NonNegative(),
	}
}

func (Budget) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("budgets").Unique().Required(),
		edge.From("category", Category.Type).Ref("budgets").Unique().Required(),
		edge.From("country", Country.Type).Ref("budgets").Unique().Required(),
		edge.From("currency", Currency.Type).Ref("budgets").Unique().Required(),
	}
}

func (Budget) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("month", "year").
			Edges("owner", "category", "country", "currency").
			Unique(),
	}
}
