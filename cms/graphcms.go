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

type GraphResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []struct {
		Message  string        `json:message`
		Location []interface{} `json:"locations"`
	} `json:"errors"`
}

// Init initialize config
func (g *GraphCMS) Init(url interface{}, key interface{}, stage interface{}, path interface{}) {
	g.url = url
	g.key = key
	g.path = fmt.Sprintf("%v", path)
	g.stage = fmt.Sprintf("%v", stage)
}

// GetNodes from the cms
func (g *GraphCMS) GetNodes() []interface{} {
	// Get Node Types
	nodeTypes := g.GetNodeTypes()
	log.Println("nodeTypes:")
	log.Println(nodeTypes)

	var allNodes []interface{}

	// Loop through each node type and pull all content
	for index := range nodeTypes {
		nodes := g.GetAllNodesByType(nodeTypes[index])
		// @todo Aggregate
		allNodes = append(allNodes, nodes)
	}

	return allNodes
}

// GetAllNodesByType will give you all nodes for a given node type.
// GraphQL does not let you select * like SQL, so we need to take the fields
// grabbed by introspection and build a graphql query to pull all fields
// for a given node type (schema model)
func (g *GraphCMS) GetAllNodesByType(name string) map[string]interface{} {
	// Get all fields for each node type (so we can build a query)
	fields := g.GetNodeFields(name)

	log.Println("fields:")
	log.Println(fields)

	// Response field structure
	type NodeFields struct {
		Type struct {
			Name   string `json:"name"`
			Kind   string `json:"kind"`
			Fields []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Type        struct {
					Name string `json:"name"`
					Kind string `json:"kind"`
				} `json:"type"`
			}
		} `json:"__type"`
	}

	// Build Query
	// Loop through all fields to build a query
	var query, fieldsQuery string
	var nodeFields NodeFields
	byteData, _ := json.Marshal(fields)
	err := json.Unmarshal(byteData, &nodeFields)

	// Handle any returned errors
	if err != nil {
		log.Fatal("Error processing node field response: \n%v", err)
	}

	query = `{
		%v {
			%v
		}
	}`
	for index := range nodeFields.Type.Fields {
		// If an object, get id of that object, otherwise just grab the field name
		if nodeFields.Type.Fields[index].Type.Kind == "OBJECT" ||
			nodeFields.Type.Fields[index].Name == "documentInStages" {
			fieldsQuery += fmt.Sprintf(`%v { id } `, nodeFields.Type.Fields[index].Name)
		} else {
			fieldsQuery += nodeFields.Type.Fields[index].Name + " "
		}
	}

	query = fmt.Sprintf(query, g.Pluralize(name), fieldsQuery)

	// Query all content for each node type
	log.Println("query: ")
	log.Println(query)

	allNodes, err := g.CallGraphApi(query, "{}")

	if err != nil {
		log.Fatalf("Error pulling all %s nodes from API: %v", name, err)
	}

	return allNodes.Data
}

// Pluralize will transform the schema model (node type) name to plural
func (g *GraphCMS) Pluralize(name string) string {
	return fmt.Sprintf("%ss", strings.ToLower(name))
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
		  }
		}
	}`
	var requestVars string = `{}`
	nodeTypes, err := g.CallGraphApi(requestQuery, requestVars)

	// Handle any returned errors
	if err != nil {
		log.Fatal("Error getting Node types from GraphCMS API: \n%v", err)
	}

	type NodeTypesResponse struct {
		Type struct {
			PossibleTypes []struct {
				Name string `json:"name"`
			} `json:"possibleTypes"`
		} `json:"__type"`
	}

	var nodeTypesResp NodeTypesResponse
	byteData, _ := json.Marshal(nodeTypes.Data)
	err = json.Unmarshal(byteData, &nodeTypesResp)

	// Handle any returned errors
	if err != nil {
		log.Fatal("Error parsing node types from response: \n%v", err)
	}

	var allTypes []string
	for index := range nodeTypesResp.Type.PossibleTypes {
		allTypes = append(allTypes, nodeTypesResp.Type.PossibleTypes[index].Name)
	}

	return allTypes
}

// Get a single schema as json
func (g *GraphCMS) GetNodeFields(name string) map[string]interface{} {
	var query = `query Type ($nodetype: String!) {
		__type(name: $nodetype) {
		  kind
		  name
		  fields {
			name
			description
			type {
			  name
			  kind
			}
		  }
		}
	}`
	var queryVars = fmt.Sprintf(`{"nodetype":"%s"}`, name)

	nodeFields, err := g.CallGraphApi(query, queryVars)

	if err != nil {
		log.Fatal("Error getting Node fields from GraphCMS API: \n%v", err)
	}

	return nodeFields.Data
}

// GetSchema returns a single schema as json
func (g *GraphCMS) GetSchema(name string) (map[string]interface{}, error) {
	var schema string
	var err error
	query, queryVars := g.GetSchemaQuery(name)

	response, err := g.CallGraphApi(query, queryVars)
	if err != nil {
		fmt.Printf("API Failre: %v", err)
	}

	// Pretty print to screen (debug)
	buff, err := json.MarshalIndent(response.Data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	schema = fmt.Sprintf("%s", buff)
	fmt.Printf("Schema: %v", schema)

	return response.Data, err
}

// GetSchemaQuery returns a GraphQL query
func (g *GraphCMS) GetSchemaQuery(name string) (string, string) {
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
	var queryVars = fmt.Sprintf(`{"nodetype":"%s"}`, name)

	fmt.Println(query)
	fmt.Println(queryVars)

	return query, queryVars
}

// GetSchemas retrives the schema and returns json
func (g *GraphCMS) GetSchemas() []interface{} {
	nodeTypes := g.GetNodeTypes()

	log.Println("GetSchemas: nodeTypes")
	log.Println(nodeTypes)

	var schemas []interface{}

	for index, _ := range nodeTypes {
		fmt.Println(index)
		fmt.Println(nodeTypes[index])
		schema, err := g.GetSchema(nodeTypes[index])
		if err != nil {
			fmt.Printf("error getting %v schema", nodeTypes[index])
		} else {
			schemas = append(schemas, schema)
		}
	}

	return schemas
}

// DownloadSchemas into a file
func (g *GraphCMS) DownloadSchemas() error {
	var err error
	schemas := g.GetSchemas()

	// Write nodes to file
	g.FileInit(g.path, fmt.Sprintf("/schemas/%v/models", g.stage), "0001.json")
	g.WriteFileJson(schemas)

	return err
}

// GetEnumerationNames from the API using introspection
func (g *GraphCMS) GetEnumerationNames() map[string]interface{} {
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

	return nodeTypes.Data
}

// Get enumerations
func (g *GraphCMS) GetEnumerations() {
	// GetEnumerationNames()
	enums := g.GetEnumeration("Tags")
	log.Println(enums)
}

// Get enumeration names from the API using introspection
func (g *GraphCMS) GetEnumeration(name string) map[string]interface{} {
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

	return nodeTypes.Data
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

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Check status
	// fmt.Println(resp)
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
	// mapBody(body)
	// fmt.Println(apiResp)
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
	data = g.GetListsV1()

	// Write lists to file
	g.FileInit(g.path, fmt.Sprintf("/content/%v/lists", g.stage), "0001.json")
	g.WriteFileJson(data)

	/* Get relations from GraphCMS and write to file */
	data = g.GetRelationsV1()

	// Write relations to file
	g.FileInit(g.path, fmt.Sprintf("/content/%v/relations", g.stage), "0001.json")
	g.WriteFileJson(data)
}
