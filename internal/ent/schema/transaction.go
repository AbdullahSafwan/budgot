package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Transaction struct {
	ent.Schema
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("transaction_type").Values("income", "expense", "transfer"),
		field.Int64("amount"),
		field.String("description").Optional(),
		field.Time("transaction_date"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("transactions").Unique().Required(),
		edge.From("category", Category.Type).Ref("transactions").Unique().Required(),
		edge.To("linked_transaction", Transaction.Type).Unique(),
	}
}
