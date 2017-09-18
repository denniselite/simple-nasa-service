package structs

type Config struct {
	Listen    string  `yaml:"listen"`
	Db        DbConfig     `yaml:"pgsql"`
	NSManager NasaServerManager `yaml:"NASAServer"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type NasaServerManager struct {
	EndPoint   string `yaml:"endPoint"`
	APIKey     string `yaml:"API-KEY"`
	APIVersion string `yaml:"APIVersion"`
}
