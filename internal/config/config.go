package config

type Config struct {
	Env   string       `yaml:"env"`
	Http  HttpConfig   `yaml:"http"`
	Auth  RemoteConfig `yaml:"auth"`
	Admin RemoteConfig `yaml:"admin"`
}

type HttpConfig struct {
	Port int `yaml:"port"`
}

type RemoteConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
