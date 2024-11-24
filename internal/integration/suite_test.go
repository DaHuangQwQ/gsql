package integration

import (
	"github.com/DaHuangQwQ/gsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	driver string
	dsn    string

	db *gsql.DB
}

func (s *Suite) SetupSuite() {
	db, err := gsql.Open(s.driver, s.dsn)
	require.NoError(s.T(), err)
	err = db.Wait()
	require.NoError(s.T(), err)
	s.db = db
}
