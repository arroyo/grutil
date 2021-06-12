/*
Copyright © 2020 John Arroyo

cms graphcms package: API client

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

// GraphResponse GraphCMS API response
type GraphResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []struct {
		Message   string        `json:message`
		Locations []interface{} `json:"locations"`
	} `json:"errors"`
}

// CallGraphAPI to make a GraphQL API call with requestQuery & requestVars
func (g *GraphCMS) CallGraphAPI(requestQuery string, requestVars string) (GraphResponse, error) {
	var url string = fmt.Sprintf("%v", g.url)
	// Safety check for GraphQL vars, API expects {}
	if requestVars == "" {
		requestVars = "{}"
	}
	requestBody := fmt.Sprintf(`{"query":"%v","variables":%v}`, g.formatQuery(requestQuery), requestVars)
	authorization := fmt.Sprintf("Bearer %v", g.key)
	bodyIoReader := strings.NewReader(requestBody)
	// fmt.Println(requestBody)

	req, err := http.NewRequest("POST", url, bodyIoReader)
	if err != nil {
		log.Fatal("error creating API request. ", err)
	}

	req.Header.Add("Authorization", authorization)
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

	// Check status
	if resp.StatusCode != 200 {
		log.Fatalf("GraphCMS server error: %v", resp.Status)
		// Write debug data to log output
		if g.Debug {
			log.Printf("Failed API query: %v", requestQuery)
			log.Printf("Failed API response body: %v", body)
		}
	}

	// Process the response
	var apiResp GraphResponse
	err = json.Unmarshal([]byte(body), &apiResp)
	if len(apiResp.Errors) > 0 {
		if len(apiResp.Data) > 0 {
			// data is sometimes returned even if there is an error. GraphCMS edge case / quirk.
			if g.Debug {
				log.Printf("API query: %v", requestQuery)
				log.Printf("API response: %v", apiResp)
			}
		} else {
			log.Fatalf("GraphCMS API returned an error: %v", apiResp.Errors[0].Message)
		}
	}

	// Write debug data to log output
	if g.Debug {
		log.Printf("API query: %v", requestQuery)
		log.Printf("API response: %v", apiResp)
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
