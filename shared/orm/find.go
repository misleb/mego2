package orm

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type FindModel struct {
	*selec[FindModel]
}

func Find(row Model) *FindModel {
	s := &FindModel{
		&selec[FindModel]{dbCommon: dbCommon{
			table:      row.TableName(),
			primaryKey: row.PrimaryKey(),
			foreignKey: row.ForeignKey(),
			model:      row,
		}},
	}
	s.columns = []string{}

	for _, v := range row.Mapping() {
		if v.NoSelect {
			continue
		}
		s.columns = append(s.columns, s.table+"."+v.Column)
	}
	s.setT(s)
	return s
}

func (d *FindModel) Query(ctx context.Context, db *sqlx.DB) error {
	if d.err != nil {
		return d.err
	}
	sqlText := d.SQL()
	rows, err := db.NamedQueryContext(ctx, sqlText, d.using)
	if err != nil {
		return err
	}
	defer rows.Close()

	// We only expect one row.
	if rows.Next() {
		err = rows.StructScan(d.model)
		if err != nil {
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}
