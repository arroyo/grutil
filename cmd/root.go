/*
Copyright © 2020 John Arroyo
*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, path, developer, directory, schema string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmsutil",
	Short: "CMS Utility",
	Long:  `A headless CMS utility for interacting with a CMS API for simple tasks like download and backup`,
}

// Config structure
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "Specify directory to save file")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Change path pointing to where config file is stored. (default $HOME/.cmsutil/config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".helloworld" (without extension).
		viper.AddConfigPath(home + "/.cmsutil")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// viper.SetEnvPrefix("CMSUTIL")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.SetDefault("backups.schemapath", home+"/.cmsutil/backups/schema")
	viper.SetDefault("backups.contentpath", home+"/.cmsutil/backups/content")
	// viper.Set("cms.host", "set override in code") // Example override

	// Validate API URL
	apiURL := viper.Get("CMS_API_URL")
	matched, err := regexp.MatchString(`^http[s]?:\/\/`, fmt.Sprintf("%v", apiURL))
	if !matched {
		log.Fatalln("Config setting CMS_API_URL does not contain a valid URL.")
	}
	if err != nil {
		log.Fatalln(err)
	}

	/*
		Fold viper config into the Config struct
		@note don't have a good way to pass this structure into a command's func
		It is not currently being used.
	*/
	var configuration Config
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		log.Fatalln(err)
	}
}
