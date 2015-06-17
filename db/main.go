package db

import (
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // allow postgres sql connections
	"github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
)

var ErrNoResults = errors.New("No results")
var ErrDestinationNotPointer = errors.New("Provided destination is not a pointer")

type Query interface {
	Get(context.Context) ([]interface{}, error)
	IsComplete(context.Context, int) bool
}

type Query2 interface {
	Select(context.Context, interface{}) error
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

// Select runs the provided query, setting all found results on dest.
func Select(ctx context.Context, query Query2, dest interface{}) error {
	dvp := reflect.ValueOf(dest)
	dv := reflect.Indirect(dvp)
	// create an intermediary slice of the correct type
	rvp := reflect.New(dv.Type())
	rv := reflect.Indirect(rvp)

	err := query.Select(ctx, rvp.Interface())

	if err != nil {
		return err
	}

	dv.Set(rv)
	return nil
}

// MustSelect is like Select, but panics on error
func MustSelect(ctx context.Context, query Query2, dest interface{}) {
	err := Select(ctx, query, dest)

	if err != nil {
		panic(err)
	}
}

// Get runs the provided query, returning the first result found, if any.
func Get(ctx context.Context, query Query2, dest interface{}) error {
	// TODO: dest must be a pointer
	// get the pointed to value
	dvp := reflect.ValueOf(dest)
	dv := reflect.Indirect(dvp)

	// create a slice of the same type as dest
	sliceType := reflect.SliceOf(dv.Type())
	rvp := reflect.New(sliceType)
	rv := reflect.Indirect(rvp)

	err := query.Select(ctx, rvp.Interface())
	if err != nil {
		return err
	}

	if rv.Len() == 0 {
		return ErrNoResults
	}

	// set the first result to the destination
	dv.Set(rv.Index(0))
	return nil
}

// MustGet is like Get, but panics on error
func MustGet(ctx context.Context, query Query2, dest interface{}) {
	err := Get(ctx, query, dest)

	if err != nil {
		panic(err)
	}
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

func setOn(src interface{}, dest interface{}) error {
	// TODO: get more rigorous with confirming dest and src are correct
	// i.e. check for settability, investigate other ways this could fail as well

	sp := reflect.ValueOf(src)
	dvp := reflect.ValueOf(dest)

	if dvp.Kind() != reflect.Ptr {
		return ErrDestinationNotPointer
	}

	dv := reflect.Indirect(dvp)
	dv.Set(sp)

	return nil
}
