package db

import (
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestHistoryAccountByAddressQuery(t *testing.T) {
	test.LoadScenario("base")
	ctx := test.Context()
	db := OpenTestDatabase()
	defer db.Close()

	Convey("AccountByAddress", t, func() {
		var account HistoryAccountRecord

		Convey("Existing record behavior", func() {
			address := "gspbxqXqEUZkiCCEFFCN9Vu4FLucdjLLdLcsV6E82Qc1T7ehsTC"
			q := HistoryAccountByAddressQuery{
				SqlQuery{db},
				address,
			}
			err := Get(ctx, q, &account)
			So(err, ShouldBeNil)
			So(account.Id, ShouldEqual, 0)
			So(account.Address, ShouldEqual, address)
		})

		Convey("Missing record behavior", func() {
			address := "not real"
			q := HistoryAccountByAddressQuery{
				SqlQuery{db},
				address,
			}
			err := Get(ctx, q, &account)
			So(err, ShouldEqual, ErrNoResults)
		})

	})
}
