package config

import (
	"embed"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed default.yml
var defaultConfigFS embed.FS

type Config struct {
	Admin  AdminConfig  `yaml:"admin"`
	Cmd    CmdConfig    `yaml:"cmd"`
	Tunnel TunnelConfig `yaml:"tunnel"`
	Acme   AcmeConfig   `yaml:"acme"`
	System SystemConfig `yaml:"system"`
}

var ConfigData Config
var ConfigFilePath string

func LoadConfig(path string) error {
	ConfigFilePath = path

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := generateDefaultConfig(path); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &ConfigData); err != nil {
		return err
	}

	return nil
}

func generateDefaultConfig(path string) error {
	data, err := defaultConfigFS.ReadFile("default.yml")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SaveConfig() error {
	data, err := yaml.Marshal(&ConfigData)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFilePath, data, 0644)
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

type CmdConfig struct {
	Port int `yaml:"port"`
}

type TunnelConfig struct {
	IP         string `yaml:"ip"`
	Port       int    `yaml:"port"`
	OpenDomain bool   `yaml:"open-domain"`
	HttpPort   int    `yaml:"http-port"`
	HttpsPort  int    `yaml:"https-port"`
}

type AcmeConfig struct {
	Email    string `yaml:"email"`
	HttpPort string `yaml:"http-port"`
}

type SystemConfig struct {
	SiteTitle      string     `yaml:"site-title"`
	OpenRegister   bool       `yaml:"open-register"`
	RegisterReview bool       `yaml:"register-review"`
	Smtp           SmtpConfig `yaml:"smtp"`
}

type SmtpConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	From      string `yaml:"from"`
	FromName  string `yaml:"from-name"`
	EnableSSL bool   `yaml:"enable-ssl"`
}
