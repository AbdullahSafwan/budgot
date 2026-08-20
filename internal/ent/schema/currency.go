package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Currency struct {
	ent.Schema
}

func (Currency) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.String("symbol").Optional(),
		field.Int("decimal_places").Default(2),
	}
}

func (Currency) Edges() []ent.Edge {
	return nil
}
