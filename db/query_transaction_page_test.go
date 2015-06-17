package db

import (
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestTransactionPageQuery(t *testing.T) {
	test.LoadScenario("base")
	ctx := test.Context()

	db := OpenTestDatabase()
	defer db.Close()

	Convey("TransactionPageQuery", t, func() {
		var records []TransactionRecord

		makeQuery := func(c string, o string, l int32) TransactionPageQuery {
			pq, err := NewPageQuery(c, o, l)

			So(err, ShouldBeNil)

			return TransactionPageQuery{
				SqlQuery:  SqlQuery{db},
				PageQuery: pq,
			}
		}

		Convey("orders properly", func() {
			// asc orders ascending by id
			MustSelect(ctx, makeQuery("", "asc", 0), &records)

			So(records, ShouldBeOrderedAscending, func(r interface{}) int64 {
				So(r, ShouldHaveSameTypeAs, TransactionRecord{})
				return r.(TransactionRecord).Id
			})

			// desc orders descending by id
			MustSelect(ctx, makeQuery("", "desc", 0), &records)

			So(records, ShouldBeOrderedDescending, func(r interface{}) int64 {
				So(r, ShouldHaveSameTypeAs, TransactionRecord{})
				return r.(TransactionRecord).Id
			})
		})

		Convey("limits properly", func() {
			// returns number specified
			MustSelect(ctx, makeQuery("", "asc", 3), &records)
			So(len(records), ShouldEqual, 3)

			// returns all rows if limit is higher
			MustSelect(ctx, makeQuery("", "asc", 10), &records)
			So(len(records), ShouldEqual, 4)
		})

		Convey("cursor works properly", func() {
			var record TransactionRecord

			// lowest id if ordered ascending and no cursor
			MustGet(ctx, makeQuery("", "asc", 0), &record)
			So(record.Id, ShouldEqual, 12884905984)

			// highest id if ordered descending and no cursor
			MustGet(ctx, makeQuery("", "desc", 0), &record)
			So(record.Id, ShouldEqual, 17179873280)

			// starts after the cursor if ordered ascending
			MustGet(ctx, makeQuery("12884905984", "asc", 0), &record)
			So(record.Id, ShouldEqual, 12884910080)

			// starts before the cursor if ordered descending
			MustGet(ctx, makeQuery("17179873280", "desc", 0), &record)
			So(record.Id, ShouldEqual, 12884914176)
		})

		Convey("restricts to address properly", func() {
			address := "gspbxqXqEUZkiCCEFFCN9Vu4FLucdjLLdLcsV6E82Qc1T7ehsTC"
			q := makeQuery("", "asc", 0)
			q.AccountAddress = address
			MustSelect(ctx, q, &records)

			So(len(records), ShouldEqual, 3)

			for _, r := range records {
				So(r.Account, ShouldEqual, address)
			}
		})

		Convey("restricts to ledger properly", func() {
			q := makeQuery("", "asc", 0)
			q.LedgerSequence = 4
			MustSelect(ctx, q, &records)

			So(len(records), ShouldEqual, 1)

			for _, r := range records {
				So(r.LedgerSequence, ShouldEqual, q.LedgerSequence)
			}
		})
	})
}
