/*
Copyright © 2020 John Arroyo

cms graphcms package

Get and download schemas and content from GraphCMS
*/

package graphcms

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Pluralize will transform the schema model (node type) name to plural
// @note this is a dumb implementation and will need to be augmented later
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
