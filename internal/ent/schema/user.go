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
	passwordHash := field.String("passwordHash").NotEmpty().Sensitive()
	createdAt := field.Time("createdAt").Default(time.Now()).Immutable()
	updatedAt := field.Time("updatedAt").Default(time.Now()).UpdateDefault(time.Now())
	isActive := field.Bool("isActive").Default(true)

	return []ent.Field{username, email, passwordHash, createdAt, updatedAt, isActive}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
