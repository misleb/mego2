package types

import "github.com/misleb/mego2/shared/orm"

type User struct {
	ID           int
	Name         string
	CurrentToken string
	Email        string
	Password     string
}

func (m *User) TableName() string {
	return "users"
}

func (m *User) PrimaryKey() string {
	return "id"
}

func (m *User) ForeignKey() string {
	return "user_id"
}

func (m *User) Mapping() []*orm.Mapping {
	return []*orm.Mapping{
		{Column: "id"},
		{Column: "name"},
		{Column: "email"},
		{
			Column:   "password",
			NoSelect: true,
			BeforeInsert: func(value string) string {
				newValue := "crypt(" + value + ", gen_salt('bf'))"
				return newValue
			},
		},
	}
}
