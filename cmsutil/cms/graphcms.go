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
	url        interface{}
	key        interface{}
	configPath string
	structure  map[string]string
	stage      string
}

type Node struct {
	TypeName         string `json:"_typeName"`
	Id               string `json:"id"`
	UpdatedAt        string `json:"updatedAt"`
	ShortDescription string `json:"shortDescription"`
	DisplayDate      string `json:"displayDate"`
	Slug             string `json:"slug"`
	Status           string `json:"status"`
	CreatedAt        string `json:"createdAt"`
	Catgeory         string `json:"category"`
	Title            string `json:"title"`
}

type NodeResponse struct {
	Out struct {
		JsonElements []interface{} `json:"jsonElements"`
	} `json:"out"`
	Cursor interface{} `json:"cursor"`
	IsFull bool        `json:"isFull"`
}

type ApiResponse struct {
	Out struct {
		JsonElements []interface{} `json:"jsonElements"`
	} `json:"out"`
	Cursor interface{} `json:"cursor"`
	IsFull bool        `json:"isFull"`
}

func (g *GraphCMS) Init(url interface{}, key interface{}, path interface{}, stage interface{}) {
	g.url = url
	g.key = key
	g.configPath = fmt.Sprintf("%v", path)
	g.stage = fmt.Sprintf("%v", stage)
}

func (g *GraphCMS) GetNodes() ([]interface{}, error) {
	var requestBody string = `{
		"fileType": "nodes",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`
	data, err := g.CallApi(requestBody, "export")

	mapBody(data)

	var nodes NodeResponse
	err = json.Unmarshal([]byte(data), &nodes)

	// fmt.Println(reflect.TypeOf(nodes).String())
	// fmt.Println(reflect.TypeOf(nodes.Out).String())

	return nodes.Out.JsonElements, err
}

func (g *GraphCMS) GetLists() ([]interface{}, error) {
	var requestBody string = `{
		"fileType": "lists",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`

	data, err := g.CallApi(requestBody, "export")

	var lists ApiResponse
	err = json.Unmarshal([]byte(data), &lists)

	fmt.Println(reflect.TypeOf(lists).String())
	fmt.Println(reflect.TypeOf(lists.Out).String())

	return lists.Out.JsonElements, err
}

func (g *GraphCMS) GetRelations() ([]interface{}, error) {
	var requestBody string = `{
		"fileType": "relations",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`

	data, err := g.CallApi(requestBody, "export")

	var relations ApiResponse
	err = json.Unmarshal([]byte(data), &relations)

	fmt.Println(reflect.TypeOf(relations).String())
	fmt.Println(reflect.TypeOf(relations.Out).String())

	return relations.Out.JsonElements, err
}

// Just for debugging, for now at least
func mapBody(body []uint8) error {
	fmt.Println(string(body))
	fmt.Println(reflect.TypeOf(body).String())

	// Unserialize
	var bodyJson interface{}
	err := json.Unmarshal([]byte(body), &bodyJson)

	fmt.Println(bodyJson)
	fmt.Println(reflect.TypeOf(bodyJson).String())

	return err
}

func (g *GraphCMS) CallApi(requestBody string, route string) ([]uint8, error) {
	url := fmt.Sprintf("%v/%v", g.url, route)
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

	return body, err
}

func (g *GraphCMS) GetSchema() {
	fmt.Println("GetSchema")
}

func (g *GraphCMS) GetSchemas() {
	fmt.Println("GetSchemas")
}

func (g *GraphCMS) GetContent() {
	fmt.Println("GetContent")

	// nodes, err := g.GetNodes()
	g.GetLists()
	g.GetRelations()
}

func (g *GraphCMS) DownloadContent() {
	/* Get nodes from GraphCMS and write to file */
	data, err := g.GetNodes()
	if err != nil {
		log.Fatalln(err)
	}

	// Write nodes to file
	g.FileInit(g.configPath, fmt.Sprintf("/content/%v/nodes", g.stage), "0001.json")
	g.WriteFileJson(data)

	/* Get lists from GraphCMS and write to file */
	data, err = g.GetLists()
	if err != nil {
		log.Fatalln(err)
	}

	// Write lists to file
	g.FileInit(g.configPath, fmt.Sprintf("/content/%v/lists", g.stage), "0001.json")
	g.WriteFileJson(data)

	/* Get relations from GraphCMS and write to file */
	data, err = g.GetRelations()
	if err != nil {
		log.Fatalln(err)
	}

	// Write relations to file
	g.FileInit(g.configPath, fmt.Sprintf("/content/%v/relations", g.stage), "0001.json")
	g.WriteFileJson(data)
}
