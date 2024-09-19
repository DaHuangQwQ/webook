package connpool

import (
	"context"
	"database/sql"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Mysql2Mongo struct {
	db      gorm.ConnPool
	mdb     *mongo.Database
	pattern *atomicx.Value[string]
}

func (m Mysql2Mongo) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql2Mongo) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql2Mongo) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql2Mongo) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	//TODO implement me
	panic("implement me")
}
