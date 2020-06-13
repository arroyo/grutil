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
	url       interface{}
	key       interface{}
	path      string
	structure map[string]string
	stage     string
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

type GraphResponse struct {
	Data struct {
		//Type map[string]json.RawMessage `json:"__type"`
		Type map[string]interface{} `json:"__type"`
		//Type3 json.RawMessage `json:"__type"`
		Schema map[string]interface{} `json:"__schema"`
	} `json:"data"`
	Errors []struct {
		Message  string        `json:message`
		Location []interface{} `json:"locations"`
	} `json:"errors"`
}

func (g *GraphCMS) Init(url interface{}, key interface{}, stage interface{}, path string) {
	g.url = url
	g.key = key
	g.path = path
	g.stage = fmt.Sprintf("%v", stage)
}

func (g *GraphCMS) GetNodes() []interface{} {
	// Get Node Types
	types := g.GetNodeTypes()
	log.Println("types:")
	log.Println(types)

	// Get all fields for each node type

	// Query all content for each node type

	// Aggregate and send back for saving

	// Old implemenation

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

// GetLists will pull the enumerations from GraphCMS
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
		log.Fatal("Error getting relations from api: \n%v", err)
	}

	return relations.Out.JsonElements
}

// Just for debugging the API response
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

// GetNodeTypes will get node types from the API using introspection
func (g *GraphCMS) GetNodeTypes() []string {
	var requestQuery string = `{
		__type(name: "Node") {
		  possibleTypes {
			name
			description
		  }
		}
	}`
	var requestVars string = `{}`
	nodeTypes, err := g.CallGraphApi(requestQuery, requestVars)

	log.Println(nodeTypes)

	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting Node types from GraphCMS API: \n%v", err)
	}

	var types []string

	fmt.Println("possibleTypes:")

	foo := nodeTypes.Data.Type["possibleTypes"]
	fooMap := foo.([]interface{})

	for _, v := range fooMap {
		vMap := v.(map[string]interface{})
		fmt.Println(vMap["name"])
		types = append(types, fmt.Sprintf("%v", vMap["name"]))
	}

	return types
}

// Get a single schema as json
func (g *GraphCMS) GetSchema(name string) string {
	var query = `query ($nodetype: String!) {
		__type(name: $nodetype) {
		  name
		  kind
		  fields {
			name
			description
			type {
			  name
			  kind
			  description
			}
		  }
		}
	  }`
	var queryVars = `{"nodetype":"Blog"}`

	fmt.Println(query)
	fmt.Println(queryVars)

	return query
}

// Get all schemas as json
func (g *GraphCMS) GetSchemas() string {
	log.Println("GetSchemas")
	nodeTypes := g.GetNodeTypes()

	log.Println("GetSchemas: nodeTypes")
	log.Println(nodeTypes)

	// Unserialize
	//var bodyJson interface{}
	//err := json.Unmarshal(nodeTypes, &bodyJson)

	//for index, _ := range nodeTypes["possibleTypes"] {
	//	fmt.Println(index)
	//	fmt.Println(nodeTypes["possibleTypes"][index]["name"])
	// Call g.GetSchema(schema) for each
	//}

	// Return array of json schemas as json as text
	var jsonSchemas = `[{
		"name": "schema",
		"complete": "me"
	}]`

	return jsonSchemas
}

// Get enumerations from the API using introspection
func (g *GraphCMS) GetEnumerationNames() interface{} {
	var requestBody string = `query Schema {
		__type(name: "Node") {
		  kind
		  name
		  possibleTypes {
			  name
			}
		}
	  }`
	nodeTypes, err := g.CallGraphApi(requestBody, "{}")

	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting enumerations from api: \n%v", err)
	}

	return nodeTypes.Data.Type
}

// Get enumerations
func (g *GraphCMS) GetEnumerations() {
	// GetEnumerationNames()
	enums := g.GetEnumeration("Tags")
	log.Println(enums)
}

// Get enumeration names from the API using introspection
func (g *GraphCMS) GetEnumeration(name string) interface{} {
	var requestBody string = `query EnumerationValues {
		__type(name: "%v") {
		  kind
		  name
		  description
		  enumValues {
			name
			description
		  }
		}
	  }`
	nodeTypes, err := g.CallGraphApi(requestBody, "{}")

	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting enumerations from api: \n%v", err)
	}

	return nodeTypes.Data.Type
}

// Prep GraphQL query
func (g *GraphCMS) formatQuery(query string) string {
	// replace new lines with a space
	query = strings.ReplaceAll(query, "\n", " ")
	// escape any double quotes
	query = strings.ReplaceAll(query, "\"", "\\\"")
	// Remove any tabs
	query = strings.ReplaceAll(query, "\t", "")

	return query
}

// Make a GraphQL API call requestQuery, requestVars
func (g *GraphCMS) CallGraphApi(requestQuery string, requestVars string) (GraphResponse, error) {
	var url string = fmt.Sprintf("%v", g.url)
	requestBody := fmt.Sprintf(`{"query":"%v","variables":%v}`, g.formatQuery(requestQuery), requestVars)
	// authorization := fmt.Sprintf("Bearer %v", g.key)
	bodyIoReader := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, bodyIoReader)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	fmt.Println(requestBody)

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Check status
	fmt.Println(resp)
	if resp.StatusCode != 200 {
		log.Fatalf("GraphCMS server error: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Process the response
	var apiResp GraphResponse
	err = json.Unmarshal([]byte(body), &apiResp)

	// Debug
	mapBody(body)
	fmt.Println(apiResp)
	//fmt.Println(apiResp.Errors)
	//fmt.Println(len(apiResp.Errors))

	if len(apiResp.Errors) > 0 {
		log.Fatalf("GraphCMS API returned an error: %v", apiResp.Errors[0].Message)
	}

	return apiResp, err
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
