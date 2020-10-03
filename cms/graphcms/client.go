/*
Copyright © 2020 John Arroyo

cms graphcms package

Connect to the GraphCMS API, graphcms.com
*/

package graphcms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// CallGraphAPI to make a GraphQL API call with requestQuery & requestVars
func (g *GraphCMS) CallGraphAPI(requestQuery string, requestVars string) (GraphResponse, error) {
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
