package db

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stellar/go-horizon/test"
)

func TestLedgerStateQuery(t *testing.T) {
	test.LoadScenario("base")
	ctx := test.Context()
	horizon := OpenTestDatabase()
	defer horizon.Close()
	core := OpenStellarCoreTestDatabase()
	defer core.Close()

	Convey("LedgerStateQuery", t, func() {
		var ls LedgerState

		q := LedgerStateQuery{
			SqlQuery{horizon},
			SqlQuery{core},
		}

		err := Get(ctx, q, &ls)
		So(err, ShouldBeNil)
		So(ls.HorizonSequence, ShouldEqual, 4)
		So(ls.StellarCoreSequence, ShouldEqual, 4)
	})
}
