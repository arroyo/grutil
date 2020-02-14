package config

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Version string
	Cms     struct {
		Provider   string
		Host       string
		PrivateKey string
	}
	Backups struct {
		SchemaPath  string
		ContentPath string
		Schemas     []string
	}
	API_URL string
	API_KEY string
}

func Load() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CMSUTIL")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	// Set default backup folders
	home, err := homedir.Dir()
	viper.SetDefault("backups.schemapath", home+"/cmsutil/backups/schema")
	viper.SetDefault("backups.contentpath", home+"/cmsutil/backups/content")
	viper.Set("cms.host", "set override in code")

	var configuration Config
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return configuration, err
}
