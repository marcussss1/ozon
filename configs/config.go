package configs

type Config struct {
	Server       Server       `yaml:"Server"`
	Postgres     Postgres     `yaml:"Postgres"`
	Cors         Cors         `yaml:"Cors"`
	Redis        Redis        `yaml:"Redis"`
	LinksService LinksService `yaml:"LinksService"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Postgres struct {
	DB             string `yaml:"db"`
	ConnectionToDB string `yaml:"connectionToDB"`
}

type Cors struct {
	AllowMethods     []string `yaml:"allowMethods"`
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
}

type Redis struct {
	Addr string `yaml:"addr"`
}

type LinksService struct {
	Addr string `yaml:"addr"`
}
