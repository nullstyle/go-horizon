package db

import (
	"testing"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestAccountByAddressQuery(t *testing.T) {
	test.LoadScenario("non_native_payment")
	ctx := test.Context()
	core := OpenStellarCoreTestDatabase()
	defer core.Close()
	history := OpenTestDatabase()
	defer history.Close()

	Convey("AccountByAddress", t, func() {
		var account AccountRecord

		notreal := "not_real"
		withtl := "gqdUHrgHUp8uMb74HiQvYztze2ffLhVXpPwj7gEZiJRa4jhCXQ"
		notl := "gspbxqXqEUZkiCCEFFCN9Vu4FLucdjLLdLcsV6E82Qc1T7ehsTC"

		q := AccountByAddressQuery{
			Core:    SqlQuery{core},
			History: SqlQuery{history},
			Address: withtl,
		}

		err := Get(ctx, q, &account)
		So(err, ShouldBeNil)

		So(account.Address, ShouldEqual, withtl)
		So(account.Seqnum, ShouldEqual, 12884901889)
		So(len(account.Trustlines), ShouldEqual, 1)

		q.Address = notl
		err = Get(ctx, q, &account)
		So(err, ShouldBeNil)
		So(len(account.Trustlines), ShouldEqual, 0)

		q.Address = notreal
		err = Get(ctx, q, &account)
		So(err, ShouldEqual, ErrNoResults)
	})
}
