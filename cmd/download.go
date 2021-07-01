/*
Package cmd download

Copyright © 2021 John Arroyo
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download schemas or content and assets. For additional help run: grutil download -h",
	Long: `There are two download options.
	grutil download content
	grutil download schemas`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing download option, for more help type, grutil download -h")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
}
