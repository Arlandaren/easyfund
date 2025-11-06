// Package schema содержит определения сущностей для Ent.
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Bank представляет банк с необходимыми полями.
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
