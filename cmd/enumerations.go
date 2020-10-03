/*
Copyright © 2020 John Arroyo
*/
package cmd

import (
	"github.com/arroyo/cmsutil/cms/graphcms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// backupCmd represents the backup command
var enumerationsCmd = &cobra.Command{
	Use:   "enumerations",
	Short: "Backup enumerations to a json file.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		enumerations(args)
	},
}

func init() {
	backupCmd.AddCommand(enumerationsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Get the schema and save it to disk
func enumerations(args []string) {
	var gcms graphcms.GraphCMS
	gcms.Init(viper.Get("CMS_API_URL"), viper.Get("CMS_API_KEY"), viper.Get("backups.stage"), viper.Get("backups.path"))
	gcms.DownloadEnumerations()
	gcms.DownloadAllEnumerations()
}
