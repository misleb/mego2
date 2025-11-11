package types

import (
	"time"

	"github.com/misleb/mego2/shared/orm"
)

type Token struct {
	ID        int
	Token     string
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
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

func (m *Token) Mapping() orm.MappingSlice {
	return orm.MappingSlice{
		{Column: "id"},
		{Column: "token"},
		{Column: "user_id"},
		{Column: "created_at"},
	}
}
