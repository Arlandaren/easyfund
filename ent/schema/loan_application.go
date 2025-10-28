package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

type LoanApplication struct {
	ent.Schema
}

func (LoanApplication) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("user_id"),
		field.Float("amount"),
		field.String("purpose"),
		field.String("financial_info").Optional(),
		field.Enum("status").Values("draft", "submitted", "under_review", "approved", "rejected", "partially_approved").Default("draft"),
		field.JSON("approval_progress", map[string]string{}).Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
