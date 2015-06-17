package db

import "golang.org/x/net/context"

type CoreAccountByAddressQuery struct {
	SqlQuery
	Address string
}

func (q CoreAccountByAddressQuery) IsComplete(ctx context.Context, alreadyDelivered int) bool {
	return alreadyDelivered > 0
}

func (q CoreAccountByAddressQuery) Select(ctx context.Context, dest interface{}) error {
	sql := CoreAccountRecordSelect.Where("accountid = ?", q.Address).Limit(1)
	return q.SqlQuery.Select(ctx, sql, dest)
}
