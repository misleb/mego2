package types

import "github.com/misleb/mego2/shared/orm"

type UserColumn string

const (
	UserColID            UserColumn = "id"
	UserColName          UserColumn = "name"
	UserColEmail         UserColumn = "email"
	UserColIsNewExternal UserColumn = "is_new_external"
	UserColPassword      UserColumn = "password"
)

type User struct {
	ID            int
	Name          string
	CurrentToken  string
	Email         string
	Password      string
	SetPassword   bool
	IsNewExternal bool `db:"is_new_external"`
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

func (m *User) Mapping() orm.MappingSlice {
	return orm.MappingSlice{
		{Column: string(UserColID)},
		{Column: string(UserColName)},
		{Column: string(UserColEmail)},
		{Column: string(UserColIsNewExternal)},
		{
			Column:   string(UserColPassword),
			NoSelect: true,
			BeforeInsert: func(value string) string {
				newValue := "crypt(" + value + ", gen_salt('bf'))"
				return newValue
			},
			BeforeFind: func(value string) string {
				return "crypt(" + value + ", password)"
			},
		},
	}
}
