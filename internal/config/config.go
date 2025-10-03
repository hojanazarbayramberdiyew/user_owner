package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug bool `yaml:"is_debug"`
	Listen  struct {
		Type string `yaml:"type"`
		Port string `yaml:"port"`
	} `yaml:"listen"`
	Storage Storage   `yaml:"storage"`
	JWT     JWTConfig `yaml:"jwt"`
}

type Storage struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"db_name"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		err := cleanenv.ReadConfig("/home/user/Desktop/task2/user_owner/config.yml", instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}

	})
	return instance
}
