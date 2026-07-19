package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Admin  AdminConfig  `yaml:"admin"`
	Cmd    CmdConfig    `yaml:"cmd"`
	Tunnel TunnelConfig `yaml:"tunnel"`
	Acme   AcmeConfig   `yaml:"acme"`
	System SystemConfig `yaml:"system"`
}

var ConfigData Config
var ConfigFilePath string

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
