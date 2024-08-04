package config

type config struct {
	Mysql
	Redis
}

type Redis struct {
	Host string
}

type Mysql struct {
	DSN string
}
