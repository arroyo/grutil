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
// @note manual implementation for now, should switch to something automated
// run version.sh and paste the results below
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of the cmsutil CLI",
	Long:  `Display the current version of the cmsutil CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
{
	"version": "v0.3.5"
	"date": "2021-06-12 01:10:47 -0800"
}
		`)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
