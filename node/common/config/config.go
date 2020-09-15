package config

import "github.com/BurntSushi/toml"

type Config struct {
	Common Common
	Gate   Gate
}

type Addr struct {
	LogicAddr   string `toml:"logicAddr"`
	ClusterAddr string `toml:"clusterAddr"`
}

type Common struct {
	CenterAddr string
}

type Gate struct {
	Addr
	ExternalAddr string `toml:"externalAddr"`
}

var conf *Config

func LoadConfig(path string) (*Config, error) {
	conf = &Config{}
	_, err := toml.DecodeFile(path, conf)
	if nil != err {
		return nil, err
	}
	return conf, nil
}

func GetConfig() *Config {
	return conf
}
