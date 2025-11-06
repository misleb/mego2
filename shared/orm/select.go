package orm

import (
	"fmt"
	"strings"
)

type selec[T any] struct {
	dbCommon
	t       *T
	results []any

	where []string
	join  []Model
}

func (d *selec[T]) setT(t *T) {
	d.t = t
}

func (d *selec[T]) Where(where AnyMap) *T {
	for k, v := range where {
		num := len(d.args) + 1
		d.where = append(d.where, strings.Replace(k, "?", fmt.Sprintf("$%d", num), 1))
		d.args = append(d.args, v)
	}
	return d.t
}

func (d *selec[T]) Join(model Model) *T {
	d.join = append(d.join, model)
	return d.t
}

func (d *selec[T]) SQL() string {
	var sb strings.Builder

	sb.WriteString("SELECT " + strings.Join(d.columns, ",") + " FROM " + d.table)

	if len(d.join) > 0 {
		for _, join := range d.join {
			sb.WriteString(" LEFT JOIN " + join.TableName() + " ON " + d.table + "." + d.primaryKey + " = " + join.TableName() + "." + d.foreignKey) //
		}
	}

	if len(d.where) > 0 {
		sb.WriteString(" WHERE " + strings.Join(d.where, " AND "))
	}

	return sb.String()
}
