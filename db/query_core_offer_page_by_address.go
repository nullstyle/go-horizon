package db

import "golang.org/x/net/context"

// CoreOfferPageByAddressQuery loads a page of active offers for the given
// address.
type CoreOfferPageByAddressQuery struct {
	SqlQuery
	PageQuery
	Address string
}

func (q CoreOfferPageByAddressQuery) Select(ctx context.Context, dest interface{}) error {
	sql := CoreOfferRecordSelect.
		Where("co.accountid = ?", q.Address).
		Limit(uint64(q.Limit))

	switch q.Order {
	case "asc":
		sql = sql.Where("co.offerid > ?", q.Cursor).OrderBy("co.offerid asc")
	case "desc":
		sql = sql.Where("co.offerid < ?", q.Cursor).OrderBy("co.offerid desc")
	}

	return q.SqlQuery.Select(ctx, sql, dest)
}

func (q CoreOfferPageByAddressQuery) IsComplete(ctx context.Context, alreadyDelivered int) bool {
	return alreadyDelivered >= int(q.Limit)
}
