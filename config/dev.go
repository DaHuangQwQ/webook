//go:build !k8s

package config

var (
	Config = config{
		Mysql: Mysql{
			DSN: "root:root@tcp(localhost:3306)/dahuang",
		},
		Redis: Redis{
			Host: "localhost",
		},
	}
)
