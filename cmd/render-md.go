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
var renderMdCmd = &cobra.Command{
	Use:   "md",
	Short: "Render schema nodes to Markdown files using a template.",
	Long: `Render schema nodes to Markdown files using a template.  
		All nodes of a given schema type will be rendered as Markdown files in the specified directory.
		Default template: Hugo style Markdown
		Default directory: render/
		`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("render md called")
	},
}

func init() {
	renderCmd.AddCommand(renderMdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
	renderMdCmd.PersistentFlags().StringVarP(&schema, "schema", "s", "", "Specify schema type to render")
	renderMdCmd.PersistentFlags().StringVarP(&template, "template", "t", "", "Template used to render nodes")
	// renderMdCmd.PersistentFlags().StringVarP(&nodeID, "node-id", "n", "", "Specify node ID to render")
}
