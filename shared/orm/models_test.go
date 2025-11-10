package orm

// TestModel is a simple model for testing
type TestModel struct {
	ID       int
	Name     string
	Age      int
	Password string
}

func (m *TestModel) TableName() string {
	return "test_table"
}

func (m *TestModel) PrimaryKey() string {
	return "id"
}

func (m *TestModel) ForeignKey() string {
	return "test_id"
}

func (m *TestModel) Mapping() []*Mapping {
	return []*Mapping{
		{Column: "id"},
		{Column: "name"},
		{Column: "age"},
		{
			Column:   "password",
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

type TestJoinModel struct {
	ID   int
	Name string
	Age  int
}

func (m *TestJoinModel) TableName() string {
	return "test_join_table"
}

func (m *TestJoinModel) PrimaryKey() string {
	return "id"
}

func (m *TestJoinModel) ForeignKey() string {
	return "test_join_id"
}

func (m *TestJoinModel) Mapping() []*Mapping {
	return []*Mapping{
		{Column: "id"},
		{Column: "name"},
		{Column: "age"},
	}
}
