package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Country struct {
	ent.Schema
}

func (Country) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
	}
}

func (Country) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("accounts", Account.Type).StorageKey(edge.Column("country_id")),
		edge.To("budgets", Budget.Type).StorageKey(edge.Column("country_id")),
	}
}
