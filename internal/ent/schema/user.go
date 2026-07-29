package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	username := field.String("username").NotEmpty().Unique()
	email := field.String("email").NotEmpty().Unique()
	password_hash := field.String("password_hash").NotEmpty().Sensitive()
	created_at := field.Time("created_at").Default(time.Now).Immutable()
	updated_at := field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now)
	is_active := field.Bool("is_active").Default(true)

	return []ent.Field{username, email, password_hash, created_at, updated_at, is_active}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
