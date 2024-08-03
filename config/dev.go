//go:build !k8s

package config

var (
	Config = config{
		Mysql: Mysql{
			host: "root:123456@tcp(127.0.0.1:3306)/test",
		},
		Redis: Redis{
			host: "localhost",
		},
	}
)
