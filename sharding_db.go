package gsql

import (
	"context"
	"database/sql"
	"math/rand"
)

type MasterSlaveDB struct {
	master *DB
	slaves []*DB
}

func (m *MasterSlaveDB) getCore() core {
	return m.master.core
}

func (m *MasterSlaveDB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	idx := rand.Intn(len(m.slaves))
	return m.slaves[idx].queryContext(ctx, query, args...)
}

func (m *MasterSlaveDB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.master.execContext(ctx, query, args...)
}

type Cluster struct {
	DBs map[string]*MasterSlaveDB
}

func (c *Cluster) getCore() core {
	//TODO implement me
	panic("implement me")
}

func (c *Cluster) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Cluster) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}
