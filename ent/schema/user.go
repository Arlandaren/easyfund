package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("email").Unique(),
		field.String("password_hash"),
		field.Enum("role").Values("borrower", "bank_risk_manager", "bank_analyst", "admin"),
		field.String("phone").Optional(),
		field.String("name").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
