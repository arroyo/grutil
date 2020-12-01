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
	"strings"
	"unicode"
)

// Pluralize will transform the schema model (node type) name to plural
// Took the 80/20 rule here and tried to capture the bulk of the cases with minimal logic
// English can be an ugly language, see rules at https://www.grammarly.com/blog/plural-nouns/
func (g *GraphCMS) Pluralize(name string) string {
	var plural string

	// Manual override in the config for difficult plurals like wife, wolf, & gas
	if v, found := g.Exceptions[strings.ToLower(name)]; found {
		plural = g.lcFirst(v)

		if g.Debug {
			fmt.Printf("manual override for %s: %s", name, plural)
		}

		return plural
	}

	lastChar := name[len(name)-1:]
	nameTransform := g.lcFirst(name)

	// If word end in 'y' and convert to 'ies'
	if lastChar == "y" {
		plural = fmt.Sprintf("%sies", nameTransform[0:len(nameTransform)-1])

	// If a word ends in ‑s, ‑x, -z or ‑o add 'es'
	} else if lastChar == "s" || lastChar == "x" || lastChar == "z" || lastChar == "o" {
		plural = fmt.Sprintf("%ses", nameTransform)
	
	// For almost all other nouns, simply add 's' to pluralize
	} else {
		plural = fmt.Sprintf("%ss", g.lcFirst(name))
	}
	
	// @todo If a word ends in ‑sh or ‑ch you add 'es'
	
	return plural
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
