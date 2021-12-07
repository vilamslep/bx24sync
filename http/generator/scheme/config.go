package scheme

type GeneratorConfig struct {
	DB        DataBaseConfig `json:"DB"`
	Web       WebConfig      `json:"Web"`
	QueryPath string         `json:"QueryPath"`
}

type DataBaseConfig struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

type WebConfig struct {
	Host string `json:"Host"`
	Port int    `json:"Port"`
}
