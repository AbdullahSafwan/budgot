package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").NotEmpty().Unique(),
		field.String("email").NotEmpty().Unique(),
		field.String("password_hash").NotEmpty().Sensitive(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Bool("is_active").Default(true),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sessions", Session.Type).
			StorageKey(edge.Column("user_id")).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("accounts", Account.Type).
			StorageKey(edge.Column("user_id")).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("categories", Category.Type).
			StorageKey(edge.Column("user_id")).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("budgets", Budget.Type).
			StorageKey(edge.Column("user_id")).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
		edge.To("transactions", Transaction.Type).
			StorageKey(edge.Column("user_id")).
			Annotations(entsql.Annotation{OnDelete: entsql.Cascade}),
	}
}
