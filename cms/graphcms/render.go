/*
Copyright © 2021 John Arroyo

graphcms package: render

Render nodes against a template
*/

package graphcms

import (
	// ht "html/template"
	// "fmt"
	"encoding/json"
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
	var err error
	g.RenderData, err = g.CallGraphAPI(query, "{}")
	if err != nil {
		log.Fatalln(err)
	}

	// Get template file, read template and render
	funcMap := tt.FuncMap{
		"title":      strings.Title,
		"json":       g.RenderJson,
		"jsonPretty": g.RenderJsonPretty,
	}
	t := tt.New(g.getTemplateFilename(template))
	t.Funcs(funcMap)
	t, err = t.ParseFiles(g.getTemplate(template))

	err = t.Execute(os.Stdout, g.RenderData)
	if err != nil {
		panic(err)
	}

	// Save rendered template to file

	return
}

// RenderJson converts the incoming data to a formatted json string
func (g *GraphCMS) RenderJson(data map[string]interface{}) string {
	jsonStr, _ := json.Marshal(data)

	return string(jsonStr)
}

// RenderJsonPretty converts the incoming data to an indented json string
func (g *GraphCMS) RenderJsonPretty(data map[string]interface{}) string {
	jsonStr, _ := json.MarshalIndent(data, "", "  ")

	return string(jsonStr)
}
