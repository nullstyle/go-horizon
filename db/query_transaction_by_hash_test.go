package db

import (
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestTransactionByHashQuery(t *testing.T) {

	Convey("TransactionByHashQuery", t, func() {
		test.LoadScenario("base")
		ctx := test.Context()
		db := OpenTestDatabase()
		defer db.Close()

		var record TransactionRecord

		Convey("Existing record behavior", func() {
			hash := "da3dae3d6baef2f56d53ff9fa4ddbc6cbda1ac798f0faa7de8edac9597c1dc0c"
			q := TransactionByHashQuery{SqlQuery{db}, hash}
			err := Get(ctx, q, &record)
			So(err, ShouldBeNil)
			So(record.TransactionHash, ShouldEqual, hash)
		})

		Convey("Missing record behavior", func() {
			hash := "not_real"
			q := TransactionByHashQuery{SqlQuery{db}, hash}
			err := Get(ctx, q, &record)
			So(err, ShouldEqual, ErrNoResults)
		})
	})
}
