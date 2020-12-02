/*
Package cmd download

Copyright © 2020 John Arroyo
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download schemas or content and assets. For additional help run: cmsutil download -h",
	Long: `There are two download options.
	cmsutil download content
	cmsutil download schemas`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing download option, for more help type, cmsutil download -h")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
	downloadCmd.PersistentFlags().StringVarP(&directory, "directory", "d", "", "Override the directory in the config file")
}
