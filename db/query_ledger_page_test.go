package db

import (
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestLedgerPageQuery(t *testing.T) {
	test.LoadScenario("base")
	ctx := test.Context()
	db := OpenTestDatabase()
	defer db.Close()

	var records []LedgerRecord

	Convey("LedgerPageQuery", t, func() {
		pq, err := NewPageQuery("0", "asc", 3)
		So(err, ShouldBeNil)

		q := LedgerPageQuery{SqlQuery{db}, pq}
		err = Select(ctx, q, &records)

		So(err, ShouldBeNil)
		So(len(records), ShouldEqual, 3)
		So(records, ShouldBeOrderedAscending, func(r interface{}) int64 {
			return r.(LedgerRecord).Id
		})

		lastLedger := records[len(records)-1]
		q.Cursor = lastLedger.Id

		err = Select(ctx, q, &records)

		So(err, ShouldBeNil)
		t.Log(records)
		So(len(records), ShouldEqual, 1)
	})
}
