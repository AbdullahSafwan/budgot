package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Account struct {
	ent.Schema
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Enum("account_type").
			Values("checking", "savings", "credit", "cash", "investment").
			Default("checking"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("accounts").Unique().Required(),
		edge.From("country", Country.Type).Ref("accounts").Unique().Required(),
		edge.From("currency", Currency.Type).Ref("accounts").Unique().Required(),
	}
}
