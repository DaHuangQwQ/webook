package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callbacks struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceID string
	Help       string
	vector     *prometheus.SummaryVec
}

func (c *Callbacks) Register(db *gorm.DB) error {
	vector := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:      c.Name,
			Subsystem: c.Subsystem,
			Namespace: c.Namespace,
			Help:      c.Help,
			ConstLabels: map[string]string{
				"db_name":     db.Name(),
				"instance_id": c.InstanceID,
			},
			Objectives: map[float64]float64{
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"type", "table"})
	prometheus.MustRegister(vector)
	c.vector = vector

	// Querys
	err := db.Callback().Query().Before("*").
		Register("prometheus_query_before", c.before("query"))
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("*").
		Register("prometheus_query_after", c.after("query"))
	if err != nil {
		return err
	}

	err = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", c.before("raw"))
	if err != nil {
		return err
	}

	err = db.Callback().Query().After("*").
		Register("prometheus_raw_after", c.after("raw"))
	if err != nil {
		return err
	}

	err = db.Callback().Create().Before("*").
		Register("prometheus_create_before", c.before("create"))
	if err != nil {
		return err
	}

	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", c.after("create"))
	if err != nil {
		return err
	}

	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", c.before("update"))
	if err != nil {
		return err
	}

	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", c.after("update"))
	if err != nil {
		return err
	}

	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", c.before("delete"))
	if err != nil {
		return err
	}

	err = db.Callback().Delete().After("*").
		Register("prometheus_delete_after", c.after("delete"))
	if err != nil {
		return err
	}
	return nil
}

// before 这里就是为了保持风格统一
func (c *Callbacks) before(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()
		db.Set("start_time", start)
	}
}

func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		// 如果上面没找到，这边必然断言失败
		start, ok := val.(time.Time)
		if !ok {
			// 没必要记录，有系统问题，可以记录日志
			return
		}
		duration := time.Since(start)
		c.vector.WithLabelValues(typ, db.Statement.Table).
			Observe(float64(duration.Milliseconds()))
	}
}
