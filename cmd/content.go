/*
Package cmd download content

Copyright © 2021 John Arroyo
*/
package cmd

import (
	"github.com/arroyo/cmsutil/cms/graphcms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backupCmd represents the backup command
var contentCmd = &cobra.Command{
	Use:   "content",
	Short: "Download all content, including assets",
	Long:  `Download all content and assets, grouped by schema type.`,
	Run: func(cmd *cobra.Command, args []string) {
		var gcms graphcms.GraphCMS
		gcms.Init(viper.Get("CMS_API_URL"), viper.Get("CMS_API_KEY"), viper.Get("backups.stage"), viper.Get("backups.path"))
		gcms.SetFlags(debug, verbose)
		gcms.DownloadContent()
	},
}

func init() {
	downloadCmd.AddCommand(contentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
