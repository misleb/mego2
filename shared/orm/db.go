package orm

type dbCommon struct {
	err        error
	args       []any
	columns    []string
	table      string
	primaryKey string
	foreignKey string
}

type AnyMap map[string]any
type Mapping struct {
	Column   string // Column name in the database
	Result   any    // query result (pointer)
	Value    any    // insert, update value
	NoSelect bool   // Hide column from SELECT queries (used for passwords and tokens)

	BeforeInsert func(value string, arg any) (string, any) // Modify value before inserting
}

type Model interface {
	Mapping() []*Mapping
	TableName() string
	PrimaryKey() string
	ForeignKey() string
}
