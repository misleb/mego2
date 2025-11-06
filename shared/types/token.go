package types

import (
	"time"

	"github.com/misleb/mego2/shared/orm"
)

type Token struct {
	ID        int
	Token     string
	UserID    int
	CreatedAt time.Time
}

func (m *Token) TableName() string {
	return "tokens"
}

func (m *Token) PrimaryKey() string {
	return "id"
}

func (m *Token) ForeignKey() string {
	return "token_id"
}

func (m *Token) Mapping() []*orm.Mapping {
	return []*orm.Mapping{
		{Column: "id", Result: &m.ID, Value: m.ID},
		{Column: "token", Result: &m.Token, Value: m.Token},
		{Column: "user_id", Result: &m.UserID, Value: m.UserID},
		{Column: "created_at", Result: &m.CreatedAt, Value: m.CreatedAt},
	}
}
