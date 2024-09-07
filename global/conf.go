package global

import "time"

type Config struct {
	Log         Log         `mapstructure:"log" json:"log" yaml:"log"`
	Seckey      Seckey      `mapstructure:"seckey" json:"seckey" yaml:"seckey"`
	CheckHeader CheckHeader `mapstructure:"checkHeader" json:"checkHeader" yaml:"checkHeader"`
	Server      Server      `mapstructure:"server" json:"server" yaml:"server"`
	Tls         Tls         `mapstructure:"tls" json:"tls" yaml:"tls"`
	Mysql       Mysql       `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis       Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
	Header      Header      `mapstructure:"header" json:"header" yaml:"header"`
}

type Tls struct {
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Cert string `mapstructure:"cert" json:"cert" yaml:"cert"`
	Key  string `mapstructure:"key" json:"key" yaml:"key"`
}

type Log struct {
	Format  string `mapstructure:"format" json:"format" yaml:"format"`
	Console bool   `mapstructure:"console" json:"console" yaml:"console"`
	Path    string `mapstructure:"path" json:"path" yaml:"path"`
	Level   string `mapstructure:"level" json:"level" yaml:"level"`
	Days    int    `mapstructure:"days" json:"days" yaml:"days"`
}

type Seckey struct {
	JwtKey     string `mapstructure:"jwtKey" json:"jwtKey" yaml:"jwtKey"`
	Jwtttl     int    `mapstructure:"jwtttl" json:"jwtttl" yaml:"jwtttl"`
	Pproftoken string `mapstructure:"pproftoken" json:"pproftoken" yaml:"pproftoken"`
}

type CheckHeader struct {
	Nonce             bool    `mapstructure:"nonce" json:"nonce" yaml:"nonce"`
	NonceCacheSeconds int     `mapstructure:"nonceCacheSeconds" json:"nonceCacheSeconds" yaml:"nonceCacheSeconds"`
	Time              bool    `mapstructure:"time" json:"time" yaml:"time"`
	Seconds           float64 `mapstructure:"seconds" json:"seconds" yaml:"seconds"`
	Sign              bool    `mapstructure:"sign" json:"sign" yaml:"sign"`
	All               bool    `mapstructure:"all" json:"all" yaml:"all"`
	Key               string  `mapstructure:"key" json:"key" yaml:"key"`
}

type Server struct {
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Addr string `mapstructure:"addr" json:"addr" yaml:"addr"`
}

type Mysql struct {
	Host                  string        `mapstructure:"host" json:"host" yaml:"host"`
	Username              string        `mapstructure:"username" json:"username" yaml:"username"`
	Password              string        `mapstructure:"password" json:"password" yaml:"password"`
	MaxOpenConnections    int           `mapstructure:"maxOpenConnections" json:"maxOpenConnections" yaml:"maxOpenConnections"`
	MaxConnectionLifeTime time.Duration `mapstructure:"maxConnectionLifeTime" json:"maxConnectionLifeTime" yaml:"maxConnectionLifeTime"`
	LogLevel              int           `mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"`
	PasswordMode          string        `mapstructure:"passwordMode" json:"passwordMode" yaml:"passwordMode"`
	Database              string        `mapstructure:"database" json:"database" yaml:"database"`
	MaxIdleConnections    int           `mapstructure:"maxIdleConnections" json:"maxIdleConnections" yaml:"maxIdleConnections"`
}

type Redis struct {
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
}

type Header struct {
	Realip    string `mapstructure:"realip" json:"realip" yaml:"realip"`
	Requestid string `mapstructure:"requestid" json:"requestid" yaml:"requestid"`
}

type EnvCfg struct {
	MyName string
	MyId   string
}
