/*
Package cmd backup

Copyright © 2020 John Arroyo
*/
package cmd

import (
	"github.com/arroyo/cmsutil/cms/graphcms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Download project schemas and content archived by date",
	Long: `Download all node, list, and relation metadata as well as download all assets.  
	Everything will be saved in the path set in your config file organized by date.`,
	Run: func(cmd *cobra.Command, args []string) {
		var gcms graphcms.GraphCMS
		gcms.Init(viper.Get("CMS_API_URL"), viper.Get("CMS_API_KEY"), viper.Get("backups.stage"), viper.Get("backups.path"))
		gcms.SetDebug(debug)
		gcms.Backup()
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
	// backupCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "Override the directory in the config file")
}
