package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type LoginAttempt struct {
	ent.Schema
}

func (LoginAttempt) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Optional(),
		field.String("ip_address").NotEmpty(),
		field.Bool("success"),
		field.Time("attempted_at").Default(time.Now).Immutable(),
	}
}

func (LoginAttempt) Edges() []ent.Edge {
	return nil
}
