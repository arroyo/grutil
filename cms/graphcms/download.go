/*
Copyright © 2020 John Arroyo

cms graphcms package: download

Retrieve and persist schemas, content & assets from GraphCMS
*/

package graphcms

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/arroyo/cmsutil/storage"
)

// WriteFileJSON write json struct to a file
func (g *GraphCMS) WriteFileJSON(data map[string]interface{}, folder string, filename string) {
	var file storage.File
	file.Init(g.path, folder, filename)
	file.WriteFileJSON(data)

	return
}

// DownloadSchemas to a file.  Download both models and enumerations.
func (g *GraphCMS) DownloadSchemas() error {
	err := g.DownloadModels()
	if err != nil {
		log.Fatalf("failed to download model schemas: %v", err)
	}

	err = g.DownloadEnumerations()
	if err != nil {
		log.Fatalf("failed to download enumeration schemas: %v", err)
	}

	return err
}

// DownloadModels to a file
func (g *GraphCMS) DownloadModels() error {
	var err error
	schemas := g.GetSchemas()

	// Write nodes to file
	g.WriteFileJSON(schemas, fmt.Sprintf("/%v/schemas/models", g.stage), "nodes.json")

	return err
}

// DownloadEnumerations to a file
func (g *GraphCMS) DownloadEnumerations() error {
	var err error
	enums := g.GetEnumerations()

	// Write nodes to file
	g.WriteFileJSON(enums, fmt.Sprintf("/%v/schemas/enumerations", g.stage), "select.json")

	return err
}

// DownloadAssets loop through nodes, look for assets and download
func (g *GraphCMS) DownloadAssets(data map[string]interface{}) {
	var node AssetNode
	var file storage.File
	file.Init(g.path, g.folder, "")

	// Loop through nodes, find assets and download them
	for index := range data["nodes"].( []interface{} ) {
		byteData, _ := json.Marshal(data["nodes"].( []interface{} )[index])
		err := json.Unmarshal(byteData, &node)
		if err != nil {
			panic(err)
		}
		if node.TypeName == "Asset" {
			err = file.DownloadFile(node.URL, node.Handle)
			if err != nil {
				panic(err)
			}
		}
	}
}

// DownloadContent from the GraphCMS
func (g *GraphCMS) DownloadContent() {
	/* Get nodes from GraphCMS and write to file */
	data := g.GetNodes()

	// Write nodes to file
	g.WriteFileJSON(data, fmt.Sprintf("/%v/content/nodes", g.stage), "0001.json")

	// Download all assets into the assets folder
	g.folder = fmt.Sprintf("/%v/content/assets", g.stage)
	g.DownloadAssets( data["Asset"].(map[string]interface{}) )

	/* Get lists from GraphCMS and write to file */
	// data = g.GetListsV1()
	// g.WriteFileJSON(data, fmt.Sprintf("/%v/content/lists", g.stage), "0001.json")

	/* Get relations from GraphCMS and write to file */
	// data = g.GetRelationsV1()
	// g.WriteFileJSON(data, fmt.Sprintf("/%v/content/relations", g.stage), "0001.json")
}
