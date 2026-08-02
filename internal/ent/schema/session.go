package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable(),
		field.Time("expires_at").Immutable(),
		field.Time("last_seen").Default(time.Now),
		field.String("ip_address").NotEmpty(),
		field.String("user_agent_hash").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("sessions").Unique().Required(),
	}
}
