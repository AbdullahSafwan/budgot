package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Transaction struct {
	ent.Schema
}

func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("amount"),
		field.String("description").Optional(),
		field.Time("transaction_date"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("transactions").Unique().Required(),
		edge.From("account", Account.Type).Ref("transactions").Unique().Required(),
		edge.From("category", Category.Type).Ref("transactions").Unique().Required(),
		// Transfer links are symmetric; create/delete both sides atomically via one service method.
		edge.To("linked_transaction", Transaction.Type).Unique(),
	}
}

func (Transaction) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("owner"),
		index.Edges("account"),
		index.Edges("category"),
	}
}
