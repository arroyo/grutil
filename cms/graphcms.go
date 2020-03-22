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
	path 	   string
	structure  map[string]string
	stage      string
}

type Node struct {
	TypeName string `json:"_typeName"`
	Id       string `json:"id"`
	Handle   string `json:"handle"`
}

type ApiResponse struct {
	Out struct {
		JsonElements []interface{} `json:"jsonElements"`
	} `json:"out"`
	Cursor interface{} `json:"cursor"`
	IsFull bool        `json:"isFull"`
	Errors []string    `json:"errors"`
}

func (g *GraphCMS) Init(url interface{}, key interface{}, stage interface{}, path string) {
	g.url = url
	g.key = key
	g.path = path
	g.stage = fmt.Sprintf("%v", stage)
}

func (g *GraphCMS) GetNodes() []interface{} {
	var requestBody string = `{
		"fileType": "nodes",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`
	nodes, err := g.CallApi(requestBody, "export")

	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting nodes from api: \n%v", err)
	}

	return nodes.Out.JsonElements
}

func (g *GraphCMS) GetLists() []interface{} {
	var requestBody string = `{
		"fileType": "lists",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`
	lists, err := g.CallApi(requestBody, "export")

	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting lists from api: \n%v", err)
	}

	// fmt.Println(reflect.TypeOf(lists).String())
	// fmt.Println(reflect.TypeOf(lists.Out).String())

	return lists.Out.JsonElements
}

func (g *GraphCMS) GetRelations() []interface{} {
	var requestBody string = `{
		"fileType": "relations",
		"cursor": {
		  "table": 0,
		  "row": 0,
		  "field": 0,
		  "array": 0
		}
	  }`
	relations, err := g.CallApi(requestBody, "export")
	
	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting relations from api: \n%v", err)
	}

	return relations.Out.JsonElements
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

// Make a GraphCMS API call
func (g *GraphCMS) CallApi(requestBody string, route string) (ApiResponse, error) {
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

	// Process the response
	var apiResp ApiResponse
	err = json.Unmarshal([]byte(body), &apiResp)

	if apiResp.Errors != nil {
		log.Fatalf("GraphCMS API returned an error: %v", apiResp.Errors[0])
	}

	return apiResp, err
}

// Get a single schema as json
func (g *GraphCMS) GetSchema(name string) string {
	var jsonSchema = `{
		"name": "schema",
		"complete": "me"
	}`
	fmt.Println(jsonSchema)
	return jsonSchema
}

// Get all schemas as json
func (g *GraphCMS) GetSchemas() string {
	// Loop through schemas in the config
	// Call g.GetSchema(schema) for each
	// Return array of json schemas as json as text
	var jsonSchemas = `[{
		"name": "schema",
		"complete": "me"
	}]`

	return jsonSchemas
}

// Loop through nodes, look for assets and download
func (g *GraphCMS) DownloadAssets(data []interface{}) {
	g.Folder = "/assets"
	var node Node

	// Loop through nodes, find assets and download them
	for index, _ := range data {
		byteData, _ := json.Marshal(data[index])
		err := json.Unmarshal(byteData, &node)
		if err != nil {
			panic(err)
		}

		if node.TypeName == "Asset" {
			url := fmt.Sprintf("https://media.graphcms.com/%v", node.Handle)
			err = g.DownloadFile(url, node.Handle)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (g *GraphCMS) DownloadContent() {
	/* Get nodes from GraphCMS and write to file */
	data := g.GetNodes()

	// Write nodes to file
	g.FileInit(g.path, fmt.Sprintf("/content/%v/nodes", g.stage), "0001.json")
	g.WriteFileJson(data)

	// Download all assets into the assets folder
	g.Folder = "/assets"
	g.DownloadAssets(data)

	/* Get lists from GraphCMS and write to file */
	data = g.GetLists()

	// Write lists to file
	g.FileInit(g.path, fmt.Sprintf("/content/%v/lists", g.stage), "0001.json")
	g.WriteFileJson(data)

	/* Get relations from GraphCMS and write to file */
	data = g.GetRelations()

	// Write relations to file
	g.FileInit(g.path, fmt.Sprintf("/content/%v/relations", g.stage), "0001.json")
	g.WriteFileJson(data)
}
