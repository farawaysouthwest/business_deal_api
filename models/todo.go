package models

import (
	"business_deal_api/graph/model"

	"github.com/gofrs/uuid"
)

type Todo struct {
	ID uuid.UUID `json:"id" db:"id"`
	Text string `json:"text" db:"text"`
	Done bool `json:"done" db:"done"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
}

type Todos []Todo

func (t Todo) Graphql () *model.Todo {
	return &model.Todo{
		ID: t.ID.String(),
		Text: t.Text,
		Done: t.Done,
	}
}