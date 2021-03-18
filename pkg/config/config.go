package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

type Config struct {
	*viper.Viper
	Uri      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Debug    bool   `mapstructure:"debug"`
}

func LoadConfig(filename string) (*Config, error) {
	conf := &Config{Viper: viper.New()}
	var configFile string

	if filename != "" {
		configFile = filename
	} else if os.Getenv("NGAC_CONFIG") != "" {
		configFile = os.Getenv("NGAC_CONFIG")
	} else {
		return nil, fmt.Errorf("Empty string is not allowed")
	}

	conf.Set("env", os.Getenv("env"))
	conf.SetConfigFile(configFile)
	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := conf.Unmarshal(conf); err != nil {
		return nil, err
	}

	go func() {
		conf.WatchConfig()
		// https://github.com/gohugoio/hugo/blob/master/watcher/batcher.go
		// https://github.com/spf13/viper/issues/609
		// for some reason this fires twice on a Win machine, and the way some editors save files.
		conf.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Configuration has been changed...")
			// only re-read if file has been modified
			if err := conf.ReadInConfig(); err != nil {
				if err == nil {
					log.Println("Reading failed after configuration update: no data was read")
				} else {
					log.Fatalf("Reading failed after configuration update: %s \n", err.Error())
				}

				return
			} else {
				log.Println("Successfully re-read config file...")
			}

		})
	}()
	return conf, nil
}

func LoadConfigFromReader(in io.Reader) (*Config, error) {
	conf := &Config{Viper: viper.New()}
	conf.Set("env", os.Getenv("env"))
	mime, err := mimetype.DetectReader(in)
	if err != nil {
		return nil, err
	}

	conf.SetConfigType(mime.Extension())
	if err := conf.ReadConfig(in); err != nil {
		return nil, err
	}

	if err := conf.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
