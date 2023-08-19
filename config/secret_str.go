package config

var RuntimeConf = RuntimeConfig{}

type RuntimeConfig struct {
	Datasource Datasource `yaml:"Datasource"`
	Server     Server     `yaml:"Server"`
}

type Datasource struct {
	Name     string `yaml:"Name"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Driver   string `yaml:"Driver"`
}

type Server struct {
	Port           int `yaml:"Port"`
	SessionTimeOut int `yaml:"SessionTimeOut"`
}
