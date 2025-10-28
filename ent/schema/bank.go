package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Bank struct {
	ent.Schema
}

func (Bank) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name"),
		field.Float("rating").Optional(),
	}
}
