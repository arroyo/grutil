/*
Copyright © 2020 John Arroyo

cms graphcms package: utilities

Get and download schemas and content from GraphCMS
*/

package graphcms

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

// Pluralize will transform the schema model (node type) name to plural
// @note this is a dumb implementation and will need to be augmented later
func (g *GraphCMS) Pluralize(name string) string {
	return fmt.Sprintf("%ss", g.lcFirst(name))
}

// lcFirst make the first character lowercase
func (g *GraphCMS) lcFirst(str string) string {  
	for i, v := range str {  
		return string(unicode.ToLower(v)) + str[i+1:]  
	}
	return ""
}

// Just for debugging the API response
func mapBody(body []uint8) error {
	fmt.Println(string(body))
	fmt.Println(reflect.TypeOf(body).String())

	// Unserialize
	var bodyJSON interface{}
	err := json.Unmarshal([]byte(body), &bodyJSON)

	fmt.Println(bodyJSON)
	fmt.Println(reflect.TypeOf(bodyJSON).String())

	return err
}
