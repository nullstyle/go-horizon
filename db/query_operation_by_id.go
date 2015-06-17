package db

import "golang.org/x/net/context"

type OperationByIdQuery struct {
	SqlQuery
	Id int64
}

func (q OperationByIdQuery) Select(ctx context.Context, dest interface{}) error {
	sql := OperationRecordSelect.Where("id = ?", q.Id).Limit(1)

	return q.SqlQuery.Select(ctx, sql, dest)
}

func (q OperationByIdQuery) IsComplete(ctx context.Context, alreadyDelivered int) bool {
	return alreadyDelivered > 0
}
