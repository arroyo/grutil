/*
Copyright © 2021 John Arroyo

cms graphcms package: render

Render given node types
*/

package graphcms

import (
	"fmt"
	"log"
)

func (g *GraphCMS) getTemplate(templateName string) {

}

// RenderTemplate with CMS data
func (g *GraphCMS) RenderTemplate(query string, template string, filename string) {

	fmt.Printf("This is the query: %v\n", query)
	fmt.Printf("This is the template: %v\n", template)
	fmt.Printf("This is the filename: %v\n", filename)

	// Get the data
	response, err := g.CallGraphAPI(query, "{}")

	if err != nil {
		log.Fatalln(err)
	}

	// @todo check err
	fmt.Println(response.Data)

	// Get template
	// templateObj := g.getTemplate(template)

	// Save rendered template to file

	return
}
