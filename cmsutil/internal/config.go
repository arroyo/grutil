package config

import (
	"github.com/spf13/viper"
	"fmt"
)

type Config struct {
	version string
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
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	viper.SetDefault("schemaPath", "./backups/schema")
	viper.SetDefault("contentPath", "./backups/content")

	testvar := viper.Get("version")
	fmt.Println(testvar)
	testvar = viper.Get("cms")
	fmt.Println(testvar)
	testvar = viper.Get("backups")
	fmt.Println(testvar)
	testvar = viper.Get("API_URL")
	fmt.Println(testvar)

	var configuration Config
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return configuration, err
}
