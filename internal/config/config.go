package config

type Config struct {
	Env  string     `yaml:"env"`
	Http HttpConfig `yaml:"http"`
}

type HttpConfig struct {
	Port int `yaml:"port"`
}
