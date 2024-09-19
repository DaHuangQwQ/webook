package connpool

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type WriteSplit struct {
	master gorm.ConnPool
	slaves []gorm.ConnPool
}

func (w *WriteSplit) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	return w.master.(gorm.ConnPoolBeginner).BeginTx(ctx, opts)
}

func (w *WriteSplit) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return w.master.PrepareContext(ctx, query)
}

func (w *WriteSplit) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.master.ExecContext(ctx, query, args...)
}

func (w *WriteSplit) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	// slaves 负载均衡
	// 轮询 加权轮询 加权轮询 平滑的加权轮询 随机 加权随机
	// 动态判定 slaves 健康情况
	//for _, v := range w.slaves {
	//	v.QueryContext(ctx, query, args...)
	//}
	panic("implement me")
}

func (w *WriteSplit) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	//TODO implement me
	panic("implement me")
}
