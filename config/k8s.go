//go:build k8s

package config

var (
	Config = config{
		Mysql: Mysql{
			DSN: "root:root@tcp(webook-mysql:3308)/dahuang",
		},
		Redis: Redis{
			Host: "webook-redis",
		},
	}
)
