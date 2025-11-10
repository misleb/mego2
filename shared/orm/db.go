package orm

type dbCommon struct {
	err        error
	columns    []string
	table      string
	primaryKey string
	foreignKey string
	model      Model
}

type AnyMap map[string]any
type Mapping struct {
	Column       string                    // Column name in the database
	NoSelect     bool                      // Hide column from SELECT queries (used for passwords and tokens)
	BeforeInsert func(value string) string // Modify value before inserting
}

type Model interface {
	Mapping() []*Mapping
	TableName() string
	PrimaryKey() string
	ForeignKey() string
}
