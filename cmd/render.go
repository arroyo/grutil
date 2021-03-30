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
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render schema nodes against a template. For additional help run: cmsutil render -h",
	Long: `There are two render types.
	cmsutil render md
	cmsutil render audio`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing render type, for more help type, cmsutil render -h")
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
}
