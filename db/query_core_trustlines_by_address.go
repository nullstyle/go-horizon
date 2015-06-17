package db

import "golang.org/x/net/context"

type CoreTrustlinesByAddressQuery struct {
	SqlQuery
	Address string
}

func (q CoreTrustlinesByAddressQuery) Select(ctx context.Context, dest interface{}) error {
	sql := CoreTrustlineRecordSelect.Where("accountid = ?", q.Address)
	return q.SqlQuery.Select(ctx, sql, dest)
}

func (q CoreTrustlinesByAddressQuery) IsComplete(ctx context.Context, alreadyDelivered int) bool {
	// this query is not stream compatible.  If we've returned any results
	// consider the query complete
	return alreadyDelivered > 0
}
