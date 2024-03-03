package config

type Config struct {
	Http HttpConfig `yaml:"http"`
}

type HttpConfig struct {
	Port int `yaml:"port"`
}
