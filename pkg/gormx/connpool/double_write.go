package connpool

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
	"webook/pkg/logger"
)

const (
	PatternSrcOnly  = "SRC_ONLY"
	PatternDstOnly  = "DST_ONLY"
	PatternSrcFirst = "SRC_FIRST"
	PatternDstFirst = "DST_FIRST"
)

var (
	errUnknownPattern = errors.New("未知的双写错误")
)

type DoubleWritePool struct {
	src     gorm.ConnPool
	dst     gorm.ConnPool
	l       logger.LoggerV1
	pattern *atomicx.Value[string]
}

func NewDoubleWritePool(src gorm.ConnPool, dst gorm.ConnPool, l logger.LoggerV1) *DoubleWritePool {
	return &DoubleWritePool{
		src:     src,
		dst:     dst,
		l:       l,
		pattern: atomicx.NewValueOf(PatternSrcOnly),
	}
}

func (d *DoubleWritePool) UpdatePattern(str string) {
	d.pattern.Store(str)
}

// 装饰器模式

func (d *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	// 事务
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{src: src, l: d.l, pattern: pattern}, err
	case PatternSrcFirst:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			// 可以考虑回滚掉 src，而后返回 error
			d.l.Error("双写目标表开启事务失败", logger.Error(err))
		}
		return &DoubleWritePoolTx{src: src, dst: dst, l: d.l, pattern: pattern}, nil
	case PatternDstFirst:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			d.l.Error("双写源表开启事务失败", logger.Error(err))
		}
		return &DoubleWritePoolTx{src: src, dst: dst, l: d.l, pattern: pattern}, nil
	case PatternDstOnly:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{dst: dst}, err
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// implement me
	panic("implement me")
}

func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// 增删改
	switch d.pattern.Load() {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, nil
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		_, err = d.src.ExecContext(ctx, query, args...)
		if err != nil {
			// log
		}
		return res, nil
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
		//return nil, errors.New("未知的双写模式")
	}
}

func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	// 查
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
		//return nil, errors.New("未知的双写模式")
	}
}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern.Load() {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
	}
}

type DoubleWritePoolTx struct {
	src     *sql.Tx
	dst     *sql.Tx
	pattern string
	l       logger.LoggerV1
}

func (d *DoubleWritePoolTx) Commit() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Commit()
	case PatternSrcFirst:
		err := d.src.Commit()
		// 如果 src 提交失败了怎么办？
		if err != nil {
			return err
		}
		if d.dst != nil {
			err1 := d.dst.Commit()
			if err1 != nil {
				d.l.Error("目标表提交事务失败")
			}
		}
		return nil
	case PatternDstFirst:
		err := d.dst.Commit()
		if err != nil {
			return err
		}
		if d.src != nil {
			err1 := d.src.Commit()
			if err1 != nil {
				d.l.Error("源表提交事务失败")
			}
		}
		return nil
	case PatternDstOnly:
		return d.dst.Commit()
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWritePoolTx) Rollback() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Rollback()
	case PatternSrcFirst:
		err := d.src.Rollback()
		// 如果 src 提交失败了怎么办？
		if err != nil {
			return err
		}
		if d.dst != nil {
			err1 := d.dst.Rollback()
			if err1 != nil {
				d.l.Error("目标表提交事务失败")
			}
		}
		return nil
	case PatternDstFirst:
		err := d.dst.Rollback()
		if err != nil {
			return err
		}
		if d.src != nil {
			err1 := d.src.Rollback()
			if err1 != nil {
				d.l.Error("源表提交事务失败")
			}
		}
		return nil
	case PatternDstOnly:
		return d.dst.Rollback()
	default:
		return errUnknownPattern
	}
}

func (d *DoubleWritePoolTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// 这个方法没办法改写
	// 我没办法返回一个双写的  sql.Stmt
	panic("双写模式写不支持")
}

func (d *DoubleWritePoolTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err == nil && d.dst != nil {
			_, err1 := d.dst.ExecContext(ctx, query, args...)
			if err1 != nil {
				d.l.Error("双写写入 dst 失败", logger.Error(err),
					logger.String("sql", query))
			}
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		// d.src != nil 开事务开成功了
		if err == nil && d.src != nil {
			_, err1 := d.src.ExecContext(ctx, query, args...)
			if err1 != nil {
				d.l.Error("双写写入 src 失败", logger.Error(err),
					logger.String("sql", query))
			}
		}
		return res, err
	default:
		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePoolTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryContext(ctx, query, args...)
	default:

		return nil, errUnknownPattern
	}
}

func (d *DoubleWritePoolTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		// 这样你没有带上错误信息
		//return &sql.Row{}
		panic(errUnknownPattern)
	}
}
