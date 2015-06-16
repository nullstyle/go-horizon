package db

import (
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // allow postgres sql connections
	"github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
)

type Query interface {
	Get(context.Context) ([]interface{}, error)
	IsComplete(context.Context, int) bool
}

type Pageable interface {
	PagingToken() string
}

type Record interface{}

// Open the postgres database at the provided url and performing an initial
// ping to ensure we can connect to it.
func Open(url string) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

// Results runs the provided query, returning all found results
func Results(ctx context.Context, query Query) ([]interface{}, error) {
	return query.Get(ctx)
}

// First runs the provided query, returning the first result if found,
// otherwise nil
func First(ctx context.Context, query Query) (interface{}, error) {
	res, err := query.Get(ctx)

	switch {
	case err != nil:
		return nil, err
	case len(res) == 0:
		return nil, nil
	default:
		return res[0], nil
	}
}

func MustFirst(ctx context.Context, q Query) interface{} {
	result, err := First(ctx, q)

	if err != nil {
		panic(err)
	}

	return result
}

func MustResults(ctx context.Context, q Query) []interface{} {
	result, err := Results(ctx, q)

	if err != nil {
		panic(err)
	}

	return result
}

func QueryGauge() metrics.Gauge {
	return globalStreamManager.queryGauge
}

// helper method suited to confirm query validity.  checkOptions ensures
// that zero or one of the provided bools ares true, but will return an error
// if more than one clause is true.
func checkOptions(clauses ...bool) error {
	hasOneSet := false

	for _, isSet := range clauses {
		if !isSet {
			continue
		}

		if hasOneSet {
			return errors.New("Invalid options: multiple are set")
		}

		hasOneSet = true
	}

	return nil
}

// Converts a typed slice to a slice of interface{}, suitable
// for return through the Get() method of Query
func makeResult(src interface{}) []interface{} {
	srcValue := reflect.ValueOf(src)
	srcLen := srcValue.Len()
	result := make([]interface{}, srcLen)

	for i := 0; i < srcLen; i++ {
		result[i] = srcValue.Index(i).Interface()
	}
	return result
}
