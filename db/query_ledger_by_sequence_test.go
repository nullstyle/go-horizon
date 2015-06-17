package db

import (
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestLedgerBySequenceQuery(t *testing.T) {

	Convey("LedgerBySequenceQuery", t, func() {
		test.LoadScenario("base")
		ctx := test.Context()
		db := OpenTestDatabase()
		defer db.Close()
		var record LedgerRecord

		Convey("Existing record behavior", func() {
			sequence := int32(2)
			q := LedgerBySequenceQuery{
				SqlQuery{db},
				sequence,
			}
			err := Get(ctx, q, &record)
			So(err, ShouldBeNil)
			So(record.Sequence, ShouldEqual, sequence)
		})

		Convey("Missing record behavior", func() {
			sequence := int32(-1)
			q := LedgerBySequenceQuery{
				SqlQuery{db},
				sequence,
			}
			err := Get(ctx, q, &record)
			So(err, ShouldEqual, ErrNoResults)
		})
	})
}
