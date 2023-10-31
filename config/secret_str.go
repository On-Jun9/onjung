package config

var RuntimeConf = RuntimeConfig{}

type RuntimeConfig struct {
	// properties.yaml
	Server    Server    `yaml:"Server"`    // Server 설정 - properties.yaml 파일
	Variables Variables `yaml:"Variables"` // Variables 설정 - properties.yaml 파일
	// secrets.yaml
	Datasource Datasource `yaml:"Datasource"` // Datasource 설정 - secrets.yaml 파일
}

type Server struct {
	Mode           string `yaml:"Mode"`
	Port           int    `yaml:"Port"`
	SessionTimeOut int    `yaml:"SessionTimeOut"`
	DBLogLevel     int    `yaml:"DBLogLevel"`
	ServerLogLevel int    `yaml:"ServerLogLevel"`
}

type Variables struct {
}

type Datasource struct {
	Name     string `yaml:"Name"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Driver   string `yaml:"Driver"`
	SslMode  string `yaml:"SslMode"`
}
