/*
Copyright © 2020 John Arroyo

cms graphcms package

Get and download schemas and content from the GraphCMS V1 (Legacy) API
*/

package cms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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

// Get Nodes from API V1 (Legacy) implementation of GraphCMS, relies on the /export endpoint
func (g *GraphCMS) GetNodesV1() []interface{} {
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

// GetListsV1 will pull the enumerations from GraphCMS API V1
func (g *GraphCMS) GetListsV1() []interface{} {
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

// GetRelationsV1 will pull the node relationships (links) from GraphCMS
func (g *GraphCMS) GetRelationsV1() []interface{} {
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

// DownloadContentV1 will download all CMS content using the legacy API
func (g *GraphCMS) DownloadContentV1() {
	/* Get nodes from GraphCMS and write to file */
	data := g.GetNodesV1()

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
