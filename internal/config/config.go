package config

import (
	"flag"
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Addr string
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH") // to get the config path from environment variable

	if configPath == "" {
		// we can check config paths in arguments or flags now

		flags := flag.String("config", "", "path to the configuration file")
		
		flag.Parse()

		configPath = *flags;

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath);
	os.IsNotExist(err){
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
    
	if err != nil {
		log.Fatalf("can not read config file: %s", err)
	}

	return &cfg
}
