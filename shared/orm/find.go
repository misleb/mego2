package orm

import (
	"context"

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
		}},
	}
	s.results = []any{}
	s.columns = []string{}

	for _, v := range row.Mapping() {
		if v.NoSelect {
			continue
		}
		s.columns = append(s.columns, row.TableName()+"."+v.Column)
		s.results = append(s.results, v.Result)
	}
	s.setT(s)
	return s
}

func (d *FindModel) Query(ctx context.Context, db *sqlx.DB) error {
	if d.err != nil {
		return d.err
	}
	sqlText := d.SQL()
	return db.QueryRowContext(ctx, sqlText, d.args...).Scan(d.results...)
}
