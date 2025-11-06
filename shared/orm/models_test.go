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
		{Column: "id", Result: &m.ID, Value: m.ID},
		{Column: "name", Result: &m.Name, Value: m.Name},
		{Column: "age", Result: &m.Age, Value: m.Age},
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
		{Column: "id", Result: &m.ID, Value: m.ID},
		{Column: "name", Result: &m.Name, Value: m.Name},
		{Column: "age", Result: &m.Age, Value: m.Age},
	}
}
