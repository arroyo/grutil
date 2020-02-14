/*
Copyright © 2020 John Arroyo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/arroyo/cmsutil/cms"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all content and assets",
	Long:  `Download all node, list, and relation metadata.  Download assets.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Begin download of CMS content and schemas...")
		download()
	},
}

func download() {
	var gcms cms.GraphCMS
	gcms.Init(viper.Get("API_URL"), viper.Get("API_KEY"), viper.Get("backups.path"), viper.Get("backups.stage"))
	gcms.DownloadContent()
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
