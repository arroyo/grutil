/*
Package cmd download

Copyright © 2020 John Arroyo
*/

package cmd

import (
	"github.com/arroyo/cmsutil/cms/graphcms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render CMS content as file(s) using a template",
	Long:  `Render the queried content nodes against a template`,
	Run: func(cmd *cobra.Command, args []string) {
		var gcms graphcms.GraphCMS
		gcms.Init(viper.Get("CMS_API_URL"), viper.Get("CMS_API_KEY"), viper.Get("backups.stage"), viper.Get("backups.path"))
		gcms.SetDebug(debug)
		gcms.RenderTemplate(query, template, outputFilename)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
	renderCmd.Flags().StringVarP(&query, "query", "q", "", "GraphQL query to pull data from CMS")
	renderCmd.Flags().StringVarP(&template, "template", "t", "", "Template to render the content against")
	renderCmd.Flags().StringVarP(&outputFilename, "output-filename", "o", "", "Filename of the output file")

	renderCmd.MarkFlagRequired("query")
	renderCmd.MarkFlagRequired("template")
	renderCmd.MarkFlagRequired("output-filename")
}
