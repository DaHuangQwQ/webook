package config

type config struct {
	Mysql
	Redis
}

type Redis struct {
	host string
}

type Mysql struct {
	host string
}
