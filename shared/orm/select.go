package orm

import (
	"strings"
)

type selec[T any] struct {
	dbCommon
	t *T

	where []string
	join  []Model
	using Model // model to use for named queries
}

func (d *selec[T]) setT(t *T) {
	d.t = t
}

func (d *selec[T]) Where(where string) *T {
	for _, mapping := range d.model.Mapping() {
		if mapping.BeforeFind != nil {
			namedArgKey := ":" + mapping.Column
			where = strings.ReplaceAll(where, namedArgKey, mapping.BeforeFind(namedArgKey))
		}
	}
	d.where = append(d.where, where)
	d.using = d.model
	return d.t
}

// override the default model for named queries
// this is useful when you want to use a different model for named queries
// for example, when you want to use a different model for a join
// or when you want to use a different model for a where clause
func (d *selec[T]) Using(model Model) *T {
	d.using = model
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
