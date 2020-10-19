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
	Short: "Download schemas or content and assets",
	Long: `There are two download options.
	cmsutil download content
	cmsutil download schemas

	for additional help run: cmsutil download -h
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing download option, for more help type, cmsutil download -h")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
