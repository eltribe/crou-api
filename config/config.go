package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	confOnce  sync.Once
	appConfig *Config
)

type Config struct {
	Host       string         `yaml:"host"`
	Port       string         `yaml:"port"`
	Develop    bool           `yaml:"develop"`
	StaticPath string         `yaml:"staticPath"`
	StaticUrl  string         `yaml:"staticUrl"`
	Cors       CorsConfig     `yaml:"cors"`
	Log        LoggerConfig   `yaml:"log"`
	Database   DatabaseConfig `yaml:"database"`
	Redis      RedisConfig    `yaml:"redis"`
	Auth       AuthConfig     `yaml:"auth"`
}

// DatabaseConfig struct
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Type     string `yaml:"type"`
	MaxOpen  int    `yaml:"maxOpen"`
	MaxIdle  int    `yaml:"maxIdle"`
}

// RedisConfig struct
type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Db   int    `yaml:"db"`
}

// CorsConfig struct
type CorsConfig struct {
	Origins     []string `yaml:"origins"`
	Methods     []string `yaml:"methods"`
	Headers     []string `yaml:"headers"`
	Credentials bool     `yaml:"credentials"`
}

type AuthConfig struct {
	Google OAuthConfig `yaml:"google"`
	Naver  OAuthConfig `yaml:"naver"`
	JWT    JWTConfig   `yaml:"jwt"`
}

type JWTConfig struct {
	Secret                   string `yaml:"secret"`
	ExpiresHours             int64  `yaml:"expiresHours"`
	RefreshTokenExpiresHours int64  `yaml:"refreshTokenExpiresHours"`
}

type OAuthConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	Redirect     string `yaml:"redirect"`
}

// logger config
type LoggerConfig struct {
	Type     string `yaml:"type"`  // options: file, stdout
	Level    string `yaml:"level"` // debug, info, error...
	FileName string `yaml:"fileName"`
}

// LoadConfigFile config 파일(yaml)을 읽고 global struct 에 저장한다.
func LoadConfigFile(filename string) *Config {
	confOnce.Do(func() {
		viper.SetConfigType("yaml")
		viper.SetConfigFile(filename)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(err)
		}
		err = viper.Unmarshal(&appConfig)
		if err != nil {
			panic(err)
		}
		log.Printf("Config Load %+v", appConfig)
	})
	return appConfig
}
