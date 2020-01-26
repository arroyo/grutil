/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all content",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		download()
	},
}

func download() { 
	var url string = os.Getenv("API_URL")
	var requestBody string = `{"query":"query ($id: ID) {\n  blog(where: {id: $id}) {\n    status\n    updatedAt\n    createdAt\n    id\n    title\n    shortDescription\n    body\n    gallery {\n      status\n      updatedAt\n      createdAt\n      id\n      handle\n      fileName\n      height\n      width\n      size\n      mimeType\n    }\n    featuredImage {\n      status\n      updatedAt\n      createdAt\n      id\n      handle\n      fileName\n      height\n      width\n      size\n      mimeType\n    }\n    catgeory\n    metaDescription\n    metaKeywords\n    tags\n    slug\n    displayDate\n    author {\n      id\n    }\n  }\n}\n","variables":{"id":"ck5be29q8ogyf099618vjf0xp"}}`
	bodyIoReader := strings.NewReader(requestBody)

	req, err := http.NewRequest("GET", url, bodyIoReader)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))

	// Unserialize
	var bodyJson interface{}
	err = json.Unmarshal([]byte(body), &bodyJson)
	fmt.Println(bodyJson)
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
