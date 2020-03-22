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
	"fmt"
	"github.com/arroyo/cmsutil/cms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all content and assets",
	Long:  `Download all node, list, and relation metadata as well as download all assets.  
	Everything will be saved in the folder set in your config file.
	An optional subfolder name can be added, e.g. cmsutil download myfolder
	`,
	Run: func(cmd *cobra.Command, args []string) {
		download(args)
	},
}

func download(args []string) {
	// Determine base file path
	var path string
	if len(args) > 0 {
		path = fmt.Sprintf("%v/%v", viper.Get("backups.path"), args[0])
	} else {
		path = fmt.Sprintf("%v", viper.Get("backups.path"))
	}

	fmt.Printf("Your downloads will be stored at %v\n", path)
	fmt.Println("Begin download of CMS content...")

	var gcms cms.GraphCMS
	gcms.Init(viper.Get("CMS_API_URL"), viper.Get("CMS_API_KEY"), viper.Get("backups.stage"), path)
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
