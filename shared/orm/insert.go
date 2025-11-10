package orm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type InsertModel[T Model] struct {
	dbCommon
}

func Insert[T Model](model T) *InsertModel[T] {
	return &InsertModel[T]{dbCommon: dbCommon{table: model.TableName(), model: model}}
}

func (i *InsertModel[T]) Query(ctx context.Context, db *sqlx.DB) error {
	rows, err := db.NamedQueryContext(ctx, i.SQL(), i.model)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.StructScan(i.model)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("no rows inserted")
	}
}

func (i *InsertModel[T]) SQL() string {
	var sb strings.Builder

	sb.WriteString("INSERT INTO " + i.table)

	mappings := i.model.Mapping()
	columns := []string{}
	values := []string{}

	for _, mapping := range mappings {
		if mapping.Column == i.model.PrimaryKey() {
			continue
		}
		columns = append(columns, mapping.Column)
		namedArg := ":" + mapping.Column

		if mapping.BeforeInsert != nil {
			values = append(values, mapping.BeforeInsert(namedArg))
		} else {
			values = append(values, namedArg)
		}
	}

	sb.WriteString(" (" + strings.Join(columns, ",") + ") VALUES (" + strings.Join(values, ",") + ")")
	sb.WriteString(" RETURNING " + i.model.PrimaryKey())

	return sb.String()
}
