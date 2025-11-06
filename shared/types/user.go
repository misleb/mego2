package types

import "github.com/misleb/mego2/shared/orm"

type User struct {
	ID           int
	Name         string
	CurrentToken string // Not in DB table
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
		{Column: "id", Result: &m.ID, Value: m.ID},
		{Column: "name", Result: &m.Name, Value: m.Name},
		{Column: "email", Result: &m.Email, Value: m.Email},
		{
			Column:   "password",
			Result:   &m.Password,
			Value:    m.Password,
			NoSelect: true,
			BeforeInsert: func(value string, arg any) (string, any) {
				newValue := "crypt(" + value + ", gen_salt('bf'))"
				return newValue, "test"
			},
		},
	}
}
