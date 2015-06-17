package db

import "golang.org/x/net/context"

// HistoryAccountPageQuery queries for a single page of HitoryAccount objects,
// in the normal collection paging style
type HistoryAccountPageQuery struct {
	SqlQuery
	PageQuery
}

// Get executes the query, returning any results
func (q HistoryAccountPageQuery) Select(ctx context.Context, dest interface{}) error {
	sql := HistoryAccountRecordSelect.
		Limit(uint64(q.Limit))

	switch q.Order {
	case "asc":
		sql = sql.Where("ha.id > ?", q.Cursor).OrderBy("ha.id asc")
	case "desc":
		sql = sql.Where("ha.id < ?", q.Cursor).OrderBy("ha.id desc")
	}

	return q.SqlQuery.Select(ctx, sql, dest)
}

// IsComplete returns true if the query considers itself complete.  In this case,
// the query will consider itself complete when it has delivered it's
// limit
func (q HistoryAccountPageQuery) IsComplete(ctx context.Context, alreadyDelivered int) bool {
	return alreadyDelivered >= int(q.Limit)
}
