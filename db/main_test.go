package db

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDBPackage(t *testing.T) {
	// ctx := test.Context()

	Convey("db.Open", t, func() {
		// TODO
	})
}

func ExampleGet() {
	db := OpenStellarCoreTestDatabase()
	defer db.Close()

	q := CoreAccountByAddressQuery{
		SqlQuery{db},
		"gspbxqXqEUZkiCCEFFCN9Vu4FLucdjLLdLcsV6E82Qc1T7ehsTC",
	}

	var account CoreAccountRecord
	err := Get(context.Background(), q, &account)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", account.Accountid)
	// Output: gspbxqXqEUZkiCCEFFCN9Vu4FLucdjLLdLcsV6E82Qc1T7ehsTC
}
