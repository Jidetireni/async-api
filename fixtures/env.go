package fixtures

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Jidetireni/async-api.git/config"
	"github.com/Jidetireni/async-api.git/store"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

type TestEnv struct {
	conf *config.Config
	Db   *sql.DB
}

func NewTestEnv(t *testing.T) *TestEnv {

	// It loads environment variables from a .env file
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}

	os.Setenv("ENV", string(config.Env_Test))

	// It initializes the configuration
	conf, err := config.NewConfig()
	require.NoError(t, err)

	// It initializes the database connection
	db, err := store.DbInit(conf)
	require.NoError(t, err)

	return &TestEnv{
		conf: conf,
		Db:   db,
	}
}

func (te *TestEnv) SetUpDb(t *testing.T) func(t *testing.T) {

	// It creates a new migrate instance using  migration files located in the migrations
	m, err := migrate.New(
		fmt.Sprintf("file://%s/migrations", te.conf.ProjectRoot),
		te.conf.DbUrl(),
	)
	require.NoError(t, err)

	// If the migrations are already up-to-date (migrate.ErrNoChange), it ignores the error.
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	return te.TearDownDb
}

func (te *TestEnv) TearDownDb(t *testing.T) {

	// It truncates the specified tables (users, refresh_tokens, report_jobs) to clean up the database.
	_, err := te.Db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", strings.Join([]string{"users", "refresh_tokens", "report_jobs"}, ",")))
	require.NoError(t, err)
}
