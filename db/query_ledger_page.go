package db

import "golang.org/x/net/context"

type LedgerPageQuery struct {
	SqlQuery
	PageQuery
}

func (q LedgerPageQuery) Select(ctx context.Context, dest interface{}) error {
	sql := LedgerRecordSelect.
		Limit(uint64(q.Limit))

	switch q.Order {
	case "asc":
		sql = sql.Where("hl.id > ?", q.Cursor).OrderBy("hl.id asc")
	case "desc":
		sql = sql.Where("hl.id < ?", q.Cursor).OrderBy("hl.id desc")
	}

	return q.SqlQuery.Select(ctx, sql, dest)
}

func (q LedgerPageQuery) IsComplete(ctx context.Context, alreadyDelivered int) bool {
	return alreadyDelivered >= int(q.Limit)
}
