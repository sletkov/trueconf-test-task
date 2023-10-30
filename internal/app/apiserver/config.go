package apiserver

type Config struct {
	BindAddr     string `toml:"bind_addr"`
	UserStoreUrl string `toml:"bind_addr"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:     ":3333",
		UserStoreUrl: "users.json",
	}
}
