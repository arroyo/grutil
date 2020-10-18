/*
Package cmd download schemas

Copyright © 2020 John Arroyo
*/
package cmd

import (
	"fmt"
	"github.com/arroyo/cmsutil/cms/graphcms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backupCmd represents the backup command
var schemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: "Download content schemas to a json file.",
	Long:  `Download model and enumeration content schemas as json files`,
	Run: func(cmd *cobra.Command, args []string) {
		schemas(args)
	},
}

func init() {
	downloadCmd.AddCommand(schemasCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Get the schema and save it to disk
func schemas(args []string) {
	var gcms graphcms.GraphCMS
	gcms.Init(viper.Get("CMS_API_URL"), viper.Get("CMS_API_KEY"), viper.Get("backups.stage"), viper.Get("backups.path"))
	err := gcms.DownloadSchemas()

	if err != nil {
		fmt.Println("error downloading schemas")
	}
}
