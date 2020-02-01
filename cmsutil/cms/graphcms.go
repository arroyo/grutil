/*
Copyright © 2020 John Arroyo

cms graphcms package

Get and download schemas and content from GraphCMS
*/

package cms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type GraphCMS struct {
	File
	url interface {}
	key interface {}
	path interface {}
}

func (g *GraphCMS) Init(url interface {}, key interface {}, path interface {}) {
	g.url = url
	g.key = key
	g.path = path
}

func (g *GraphCMS) GetNodes() {
	var requestBody string = `{
		"fileType": "nodes",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`
	
	g.CallApi(requestBody)
}

func (g *GraphCMS) GetLists() {
	var requestBody string = `{
		"fileType": "lists",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`
	
	g.CallApi(requestBody)
}

func (g *GraphCMS) GetRelations() {
	var requestBody string = `{
		"fileType": "relations",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`
	
	g.CallApi(requestBody)
}

func (g *GraphCMS) CallApi(requestBody string) {
	url := fmt.Sprintf("%v/export", g.url)
	authorization := fmt.Sprintf("Bearer %v", g.key)
	bodyIoReader := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, bodyIoReader)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

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
	fmt.Println(reflect.TypeOf(body).String())

	// Unserialize
	var bodyJson interface{}
	err = json.Unmarshal([]byte(body), &bodyJson)

	fmt.Println(bodyJson)
	fmt.Println(reflect.TypeOf(bodyJson).String())
}

func (g *GraphCMS) GetSchema() {
	fmt.Println("GetSchema")
}

func (g *GraphCMS) GetSchemas() {
	fmt.Println("GetSchemas")	
}

func (g *GraphCMS) GetContent() {
	fmt.Println("GetContent")

	g.GetNodes()
	g.GetLists()
	g.GetRelations()
}

func (g *GraphCMS) DownloadContent() {
	g.GetContent()
	g.Folder = "/content"; g.Filename = "file.json";
	g.WriteFile(fmt.Sprintf("%v", g.path))
}
