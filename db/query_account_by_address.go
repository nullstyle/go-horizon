package db

import "golang.org/x/net/context"

// AccountByAddressQuery represents a query that retrieves a composite
// of the CoreAccount and the HistoryAccount associated with an address.
type AccountByAddressQuery struct {
	History SqlQuery
	Core    SqlQuery
	Address string
}

// IsComplete returns true when the query considers itself finished.
func (q AccountByAddressQuery) IsComplete(ctx context.Context, alreadyDelivered int) bool {
	return alreadyDelivered > 0
}

func (q AccountByAddressQuery) Select(ctx context.Context, dest interface{}) error {
	var result AccountRecord
	var cq Query2

	cq = HistoryAccountByAddressQuery{q.History, q.Address}
	err := Get(ctx, cq, &result.HistoryAccountRecord)
	if err != nil {
		return err
	}

	cq = CoreAccountByAddressQuery{q.Core, q.Address}
	err = Get(ctx, cq, &result.CoreAccountRecord)
	if err != nil {
		return err
	}

	cq = CoreTrustlinesByAddressQuery{q.Core, q.Address}
	err = Select(ctx, cq, &result.Trustlines)
	if err != nil {
		return err
	}
	setOn([]AccountRecord{result}, dest)
	return nil
}
