package config

import "github.com/spf13/viper"

type Config struct {
	cms struct {
		provider   string
		host       string
		privateKey string
	}

	backups struct {
		schemaPath  string
		contentPath string
		schemas     []string
	}
}

func InitViper() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("schemaPath", "./backups/schema")
	viper.SetDefault("contentPath", "./backups/content")

	var consts Config
	err = viper.Unmarshal(&consts)
	return consts, err
}
