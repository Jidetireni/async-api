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
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
	os.Setenv("ENV", string(config.Env_Test))
	conf, err := config.NewConfig()
	require.NoError(t, err)

	db, err := store.DbInit(conf)
	require.NoError(t, err)

	return &TestEnv{
		conf: conf,
		Db:   db,
	}
}

func (te *TestEnv) SetUpDb(t *testing.T) func(t *testing.T) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s/migrations", te.conf.ProjectRoot),
		te.conf.DbUrl(),
	)

	require.NoError(t, err)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	return te.TearDownDb
}

func (te *TestEnv) TearDownDb(t *testing.T) {
	_, err := te.Db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", strings.Join([]string{"users", "refresh_tokens", "report_jobs"}, ",")))
	require.NoError(t, err)
}
