package orm

import (
	"context"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type InsertModel[T Model] struct {
	dbCommon
	model T
}

func Insert[T Model](model T) *InsertModel[T] {
	return &InsertModel[T]{model: model, dbCommon: dbCommon{table: model.TableName()}}
}

func (i *InsertModel[T]) Query(ctx context.Context, db *sqlx.DB) error {
	// We assume that the first mapping is the primary key. FIXME: This is a hack.
	return db.QueryRowContext(ctx, i.SQL(), i.args...).Scan(i.model.Mapping()[0].Result)
}

func (i *InsertModel[T]) SQL() string {
	var sb strings.Builder

	sb.WriteString("INSERT INTO " + i.table)

	mappings := i.model.Mapping()
	i.args = []any{}
	columns := []string{}
	values := []string{}
	count := 0

	for _, mapping := range mappings {
		if mapping.Column == i.model.PrimaryKey() {
			continue
		}
		count = count + 1
		columns = append(columns, mapping.Column)

		numArg := "$" + strconv.Itoa(count)
		if mapping.BeforeInsert != nil {
			value, arg := mapping.BeforeInsert(numArg, mapping.Value)
			i.args = append(i.args, arg)
			values = append(values, value)
		} else {
			i.args = append(i.args, mapping.Value)
			values = append(values, numArg)
		}
	}

	sb.WriteString(" (" + strings.Join(columns, ",") + ") VALUES (" + strings.Join(values, ",") + ")")
	sb.WriteString(" RETURNING " + i.model.PrimaryKey())

	return sb.String()
}
