package config

type Config struct {
	Api      Api      `json:"api"`
	Postgres Postgres `json:"postgres"`
}

type Api struct {
	Port int `json:"port"`
}

type Postgres struct {
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Base     string `json:"base"`
	Sslmode  string `json:"sslmode"`
}
