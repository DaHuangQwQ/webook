//go:build k8s

package config

var (
	Config = config{
		Mysql: Mysql{
			host: "root:root@tcp(webook-mysql:3308)/dahuang",
		},
		Redis: Redis{
			host: "webook-redis",
		},
	}
)
