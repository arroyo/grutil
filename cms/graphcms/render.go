/*
Copyright © 2021 John Arroyo

graphcms package: render

Render given nodes against a template
*/

package graphcms

import (
	// ht "html/template"
	// "fmt"
	"log"
	"os"
	"regexp"
	"strings"
	tt "text/template"	
)

// getTemplateFilename gets the filename from a template file path
func (g *GraphCMS) getTemplateFilename(templateName string) string {
	re := regexp.MustCompile(`[\w-]+\.[\w]{1,4}`)
	return re.FindString(templateName)
}

// getTemplate determines which template from the filestsystem to use
// Fatal if the template does not exist
// @todo Check embedded templates
func (g *GraphCMS) getTemplate(templateName string) string {
	if _, err := os.Stat(templateName); os.IsNotExist(err) {
		log.Fatalf("template %v does not exist", templateName)
	}

	return templateName
}

// RenderTemplate with CMS data
func (g *GraphCMS) RenderTemplate(query string, template string, outputFilename string) {
	// Get the queried data from GraphCMS
	response, err := g.CallGraphAPI(query, "{}")
	if err != nil {
		log.Fatalln(err)
	}

	// Get template file, read template and render
	funcMap := tt.FuncMap {
		"title": strings.Title,
		"json": RenderJson,
		"oddOrEven": OddOrEven,
    }
	t := tt.New(g.getTemplateFilename(template))
	t.Funcs(funcMap)
	t, err = t.ParseFiles(g.getTemplate(template))
	
	err = t.Execute(os.Stdout, response)
	if err != nil {
		panic(err)
	}

	// Save rendered template to file

	return
}

func RenderJson(s string) string {
	if s == "" {
		// @todo grab the full response data
	} else {
		// @todo parse the string to determine which json segment to render
	}

	if len(s)%2 == 0 {
			return "even"
	} else {
			return "odd"
	}

}

func OddOrEven(s string) string {

	if len(s)%2 == 0 {
			return "even"
	} else {
			return "odd"
	}

}