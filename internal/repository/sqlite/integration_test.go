//go:build integration

package sqlite

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SQLiteIntegrationSuite struct {
	suite.Suite

	migrationOpts *MigrationOptions
	logger        *slog.Logger

	repository *Repository
}

func (s *SQLiteIntegrationSuite) SetupSuite() {
	var err error
	s.repository, err = NewRepository("integration_test.db",
		s.logger,
		s.migrationOpts,
	)
	s.Require().NoError(err)
}

func (s *SQLiteIntegrationSuite) SetupTest() {
	s.repository.APIKeyRepository.conn.Migrator().DropTable(tableList()...)
	s.Require().NoError(MigrateAndSeed(s.repository.APIKeyRepository.conn, s.logger, s.migrationOpts))
}

func (s *SQLiteIntegrationSuite) TearDownSuite() {
	s.NoError(os.Remove("integration_test.db"))
}

func TestSQLiteIntegration(t *testing.T) {
	sqliteIntegrationSuite := &SQLiteIntegrationSuite{
		migrationOpts: &MigrationOptions{
			APIKey: "foo",
		},
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
	// ensure always cleanup integration test db.
	defer func() {
		if r := recover(); r != nil {
			sqliteIntegrationSuite.TearDownSuite()
		}
	}()
	suite.Run(t, sqliteIntegrationSuite)
}
