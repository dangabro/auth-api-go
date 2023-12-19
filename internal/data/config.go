package data

type DbConfig struct {
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Machine  string `json:"machine"`
	Port     int32  `json:"port"`
	PoolSize int32  `json:"poolSize"`
}

type Config struct {
	Port            string   `json:"port"`
	TokenDurationMs int64    `json:"tokenDurationMs"`
	Cors            bool     `json:"cors"`
	LogLevel        string   `json:"logLevel"`
	Db              DbConfig `json:"db"`
}
