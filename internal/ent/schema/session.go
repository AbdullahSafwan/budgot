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
	id := field.String("id").NotEmpty().Immutable()
	expires_at := field.Time("expires_at").Immutable()
	last_seen := field.Time("last_seen").Default(time.Now)
	ip_address := field.String("ip_address").NotEmpty()
	user_agent_hash := field.String("user_agent_hash").NotEmpty()
	created_at := field.Time("created_at").Default(time.Now).Immutable()

	return []ent.Field{id, expires_at, last_seen, ip_address, user_agent_hash, created_at}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	edge := edge.From("owner", User.Type).Ref("sessions").Unique().Required()
	return []ent.Edge{edge}
}
