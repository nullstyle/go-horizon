// Package test contains simple test helpers that should not
// have any dependencies on horizon's packages.  think constants,
// custom matchers, generic helpers etc.
package test

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	glog "github.com/stellar/go-horizon/log"
	"golang.org/x/net/context"
)

//go:generate go get github.com/jteeuwen/go-bindata/...
//go:generate go-bindata -pkg test scenarios

const (
	// DefaultTestDatabaseURL is the default db url for horizon in a test context
	DefaultTestDatabaseURL = "postgres://localhost:5432/horizon_test?sslmode=disable"
	// DefaultTestStellarCoreDatabaseURL is the default db url for stellar-core in a test context
	DefaultTestStellarCoreDatabaseURL = "postgres://localhost:5432/stellar-core_test?sslmode=disable"
)

// DatabaseURL returns the database connection the url any test
// use when connecting to the history/horizon database
func DatabaseURL() string {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = DefaultTestDatabaseURL
	}

	return databaseURL
}

// StellarCoreDatabaseURL returns the database connection the url any test
// use when connecting to the stellar-core database
func StellarCoreDatabaseURL() string {
	databaseURL := os.Getenv("STELLAR_CORE_DATABASE_URL")

	if databaseURL == "" {
		databaseURL = DefaultTestStellarCoreDatabaseURL
	}

	return databaseURL
}

// OpenDatabase opens a database, panicing if it cannot
func OpenDatabase(dsn string) *sqlx.DB {
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		log.Panic(err)
	}

	return db
}

// LoadScenario populates the test databases with pre-created scenarios.  Each
// scenario is in the scenarios subfolder of this package and are a pair of
// sql files, one per database.
func LoadScenario(scenarioName string) {
	scenarioBasePath := "scenarios/" + scenarioName
	horizonPath := scenarioBasePath + "-horizon.sql"
	stellarCorePath := scenarioBasePath + "-core.sql"

	loadSqlFile(DatabaseURL(), horizonPath)
	loadSqlFile(StellarCoreDatabaseURL(), stellarCorePath)
}

func loadSqlFile(url string, path string) {
	sql, err := Asset(path)

	if err != nil {
		log.Panic(err)
	}

	reader := bytes.NewReader(sql)
	cmd := exec.Command("psql", url)
	cmd.Stdin = reader

	err = cmd.Run()

	if err != nil {
		log.Panic(err)
	}

}

// Context provides a context suitable for testing in tests that do not create
// a full App instance (in which case your tests should be using the app's
// context).  This context has a logger bound to it suitable for testing.
func Context() context.Context {
	return glog.Context(context.Background(), testLogger)
}

// ContextWithLogBuffer returns a context and a buffer into which the new, bound
// logger will write into.  This method allows you to inspect what data was
// logged more easily in your tests.
func ContextWithLogBuffer() (context.Context, *bytes.Buffer) {
	output := new(bytes.Buffer)
	l, _ := glog.New()
	l.Logger.Out = output
	l.Logger.Level = logrus.DebugLevel

	ctx := glog.Context(context.Background(), l)
	return ctx, output

}
