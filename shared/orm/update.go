package orm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UpdateModel[T Model] struct {
	dbCommon
	sets []string
}

func Update[T Model](model T) *UpdateModel[T] {
	return &UpdateModel[T]{dbCommon: dbCommon{table: model.TableName(), model: model}}
}

func (i *UpdateModel[T]) Query(ctx context.Context, db *sqlx.DB) error {
	if len(i.sets) == 0 {
		return fmt.Errorf("columns to set are required")
	}
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
		return fmt.Errorf("no rows updated")
	}
}

func (i *UpdateModel[T]) Set(sets []string) *UpdateModel[T] {
	i.sets = sets
	return i
}

func (i *UpdateModel[T]) SQL() string {
	var sb strings.Builder

	sb.WriteString("UPDATE " + i.table)

	// TODO: Support full model update (consitdering special fields like password)
	columns := make([]string, len(i.sets))

	for x, column := range i.sets {
		namedArg := ":" + column
		if mapping := i.model.Mapping().Find(column); mapping != nil {
			if mapping.BeforeInsert != nil {
				namedArg = mapping.BeforeInsert(namedArg)
			}
		}
		columns[x] = column + " = " + namedArg
	}

	sb.WriteString(" SET " + strings.Join(columns, ","))
	sb.WriteString(" WHERE " + i.model.PrimaryKey() + " = :" + i.model.PrimaryKey())
	sb.WriteString(" RETURNING " + i.model.PrimaryKey())

	return sb.String()
}
