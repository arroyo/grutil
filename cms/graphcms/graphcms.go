/*
Copyright © 2020 John Arroyo

cms graphcms package

Get and download schemas and content from GraphCMS
*/

package graphcms

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// GraphCMS struct
type GraphCMS struct {
	url       interface{}
	key       interface{}
	path      string
	structure map[string]string
	stage     string
}

// Node content node
type Node struct {
	TypeName string `json:"_typeName"`
	ID       string `json:"id"`
	Handle   string `json:"handle"`
}

// GraphResponse GraphCMS API response
type GraphResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []struct {
		Message  string        `json:message`
		Location []interface{} `json:"locations"`
	} `json:"errors"`
}

// Init initialize config
func (g *GraphCMS) Init(url interface{}, key interface{}, stage interface{}, path interface{}) {
	g.url = url
	g.key = key
	g.path = fmt.Sprintf("%v", path)
	g.stage = fmt.Sprintf("%v", stage)
}

// GetNodes from the cms
func (g *GraphCMS) GetNodes() []interface{} {
	// Get Node Types
	nodeTypes := g.GetNodeTypes()
	log.Println("nodeTypes:")
	log.Println(nodeTypes)

	var allNodes []interface{}

	// Loop through each node type and pull all content
	for index := range nodeTypes {
		nodes := g.GetAllNodesByType(nodeTypes[index])
		// @todo Aggregate
		allNodes = append(allNodes, nodes)
	}

	return allNodes
}

// subfieldFormat will apply special rules to handle fields that both 
// can and cannot be accounted for based on introspection.
func (g *GraphCMS) subfieldFormat(field NodeSubfield) string {
	format := ""

	if len(field.Args) > 1 {
		return field.Name + " { id }\n"
	}

	// rgba is a normalized field with no distinguisging introspection data
	// distance is based on a map input field that is irrelavant for a backup
	switch field.Name {
	case "rgba":
		format = "rgba { r g b a }\n"
	case "distance":
		format = ""
	default:
		format = field.Name + "\n"
	}

	return format
}

// GetAllNodesByType will give you all nodes for a given node type.
// GraphQL does not let you select * like SQL, so we need to take the fields
// grabbed by introspection and build a graphql query to pull all fields
// for a given node type (schema model)
func (g *GraphCMS) GetAllNodesByType(name string) map[string]interface{} {
	// Get all fields for each node type (so we can build a query)
	fields := g.GetNodeFields(name)

	// Build Query
	// Loop through all fields to build a query
	var query, fieldsQuery string
	var nodeFields NodeFields
	byteData, _ := json.Marshal(fields)
	err := json.Unmarshal(byteData, &nodeFields)

	log.Printf("Fields for schema %v:", name)
	log.Println(nodeFields)

	// Handle any returned errors
	if err != nil {
		log.Fatal("error processing node fields: \n%v", err)
	}

	// Query builder, create GraphQL query to pull all nodes and content fields
	query = `query {
		%v {
			%v
		}
	}`
	for index := range nodeFields.Type.Fields {
		// Debug
		log.Printf("name: %v, fields: %v, args: %v", nodeFields.Type.Fields[index].Name, len(nodeFields.Type.Fields[index].Type.Fields), len(nodeFields.Type.Fields[index].Args))
		
		// If an object, get id of that object, otherwise just grab the field name
		if len(nodeFields.Type.Fields[index].Type.Fields) > 0 {
			fields := ""
			for _, field := range nodeFields.Type.Fields[index].Type.Fields {
				fields += g.subfieldFormat(field)
			}
			fieldsQuery += fmt.Sprintf("%v { %v } \n", nodeFields.Type.Fields[index].Name, fields)
		} else if len(nodeFields.Type.Fields[index].Args) > 1 {
			fieldsQuery += fmt.Sprintf("%v { id } \n", nodeFields.Type.Fields[index].Name)
		} else {
			fieldsQuery += nodeFields.Type.Fields[index].Name + "\n"
		}
	}

	query = fmt.Sprintf(query, g.Pluralize(name), fieldsQuery)

	// Query all content for each node type
	log.Println("query: ")
	log.Println(query)

	allNodes, err := g.CallGraphAPI(query, "{}")

	if err != nil {
		log.Fatalf("Error pulling all %s nodes from API: %v", name, err)
	}

	return allNodes.Data
}

// GetNodeTypes will get node types from the API using introspection
func (g *GraphCMS) GetNodeTypes() []string {
	var requestQuery string = `{
		__type(name: "Node") {
		  possibleTypes {
			name
		  }
		}
	}`
	var requestVars string = `{}`
	nodeTypes, err := g.CallGraphAPI(requestQuery, requestVars)

	// Handle any returned errors
	if err != nil {
		log.Fatal("Error getting Node types from GraphCMS API: \n%v", err)
	}

	type NodeTypesResponse struct {
		Type struct {
			PossibleTypes []struct {
				Name string `json:"name"`
			} `json:"possibleTypes"`
		} `json:"__type"`
	}

	var nodeTypesResp NodeTypesResponse
	byteData, _ := json.Marshal(nodeTypes.Data)
	err = json.Unmarshal(byteData, &nodeTypesResp)

	// Handle any returned errors
	if err != nil {
		log.Fatal("Error parsing node types from response: \n%v", err)
	}

	var allTypes []string
	for index := range nodeTypesResp.Type.PossibleTypes {
		allTypes = append(allTypes, nodeTypesResp.Type.PossibleTypes[index].Name)
	}

	return allTypes
}

// NodeSubfield structure
type NodeSubfield struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Args []struct {
		Name string `json:"name"`
		Description string `json:"description"`
		DefaultValue string `json:"defaultValue"`
	} `json:"args"`
}

// NodeFields structure
type NodeFields struct {
	Type struct {
		Name   string `json:"name"`
		Kind   string `json:"kind"`
		Fields []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Args []struct {
				Name string `json:"name"`
				Description string `json:"description"`
			  	DefaultValue string `json:"defaultValue"`
			} `json:"args"`
			Type struct {
				Name string `json:"name"`
				Kind string `json:"kind"`
				Description string `json:"description"`
				Fields []NodeSubfield `json:"fields"`
			  	OfType struct {
					Name string `json:"name"`
					Description string `json:"description"`
			  	}`json:"ofType"`
			  	Interfaces []struct {
					Name string `json:"name"`
					Description string `json:"description"`
					Typename string `json:"__typename"`
			  	} `json:"interfaces"`
			} `json:"type"`
		} `json:"fields"`
	} `json:"__type"`
}

// GetNodeFields returns a single schema as json
func (g *GraphCMS) GetNodeFields(name string) map[string]interface{} {
	var query = `query GetNodeByTypeVerbose($nodetype: String!) {
		__type(name: $nodetype) {
		  name
		  kind
		  fields {
			name
			description
			args {
			  name
			  description
			  defaultValue
			}
			type {
			  name
			  kind
			  description
			  fields {
				name
				description
				args {
					name
					description
					defaultValue
				}
			  }
			  possibleTypes {
				name
				description
			  }
			  ofType {
				name
				description
			  }
			  interfaces {
				name
				description
				__typename
			  }
			}
		  }
		}
	  }`
	var queryVars = fmt.Sprintf(`{"nodetype":"%s"}`, name)
	nodeFields, err := g.CallGraphAPI(query, queryVars)

	if err != nil {
		log.Fatal("Error getting Node fields from GraphCMS API: \n%v", err)
	}

	return nodeFields.Data
}

// GetSchema returns a single schema as json
func (g *GraphCMS) GetSchema(name string) (map[string]interface{}, error) {
	var schema string
	var err error
	query, queryVars := g.GetSchemaQuery(name)

	response, err := g.CallGraphAPI(query, queryVars)
	if err != nil {
		fmt.Printf("API Failre: %v", err)
	}

	// Pretty print to screen (debug)
	buff, err := json.MarshalIndent(response.Data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	schema = fmt.Sprintf("%s", buff)
	fmt.Printf("Schema: %v", schema)

	return response.Data, err
}

// GetSchemaQuery returns a GraphQL query
func (g *GraphCMS) GetSchemaQuery(name string) (string, string) {
	var query = `query ($nodetype: String!) {
		__type(name: $nodetype) {
		  name
		  kind
		  fields {
			name
			description
			type {
			  name
			  kind
			  description
			}
		  }
		}
	  }`
	var queryVars = fmt.Sprintf(`{"nodetype":"%s"}`, name)

	fmt.Println(query)
	fmt.Println(queryVars)

	return query, queryVars
}

// GetSchemas retrives the schema and returns json
func (g *GraphCMS) GetSchemas() []interface{} {
	nodeTypes := g.GetNodeTypes()

	log.Println("GetSchemas: nodeTypes")
	log.Println(nodeTypes)

	var schemas []interface{}

	for index, _ := range nodeTypes {
		fmt.Println(index)
		fmt.Println(nodeTypes[index])
		schema, err := g.GetSchema(nodeTypes[index])
		if err != nil {
			fmt.Printf("error getting %v schema", nodeTypes[index])
		} else {
			schemas = append(schemas, schema)
		}
	}

	return schemas
}

// GetAllEnumerations from the API using introspection
// This returns numerous ENUMS including any system enumerations
func (g *GraphCMS) GetAllEnumerations() []interface{} {
	var query string = `query SchemaTypes {
		__schema {
			types {
				name
				kind
				description
				enumValues {
					name
					description
				}
		  	}
		}
	}`
	schemaTypesResp, err := g.CallGraphAPI(query, "{}")

	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting enumerations from api: \n%v", err)
	}

	type SchemaTypes struct {
		Schema struct {
			Types []struct {
				Name string `json:"name"`
				Kind string `json:"kind"`
				Description string `json:"description"`
				EnumValues []struct {
					Name string `json:"name"`
					Description string `json:"description"`
				} `json:"enumValues"`
			} `json:"types"`
		} `json:"__schema"`
	}

	var schemaTypes SchemaTypes
	byteData, _ := json.Marshal(schemaTypesResp.Data)
	err = json.Unmarshal(byteData, &schemaTypes)

	// Handle any returned errors
	if err != nil {
		log.Fatal("Error parsing schema types from response: \n%v", err)
	}

	var allTypes []interface{}
	for _, schemaType := range schemaTypes.Schema.Types {
		if schemaType.Kind == "ENUM" {
			allTypes = append(allTypes, schemaType)
		}
	}

	return allTypes
}


// GetAllEnumerations from the CMS, including system ENUMS
/*
func (g *GraphCMS) GetAllEnumerations() []interface{} {	
	names := g.GetEnumerationNames()
	var enums []interface{}

	for _, name := range names {
		enum := g.GetEnumeration(name)
		enums = append(enums, enum)
	}

	return enums
}
*/

// GetEnumerations from the CMS based on enumerations defined in your config
func (g *GraphCMS) GetEnumerations() []interface{} {
	var enums []interface{}
	enumConfig := viper.GetStringSlice("backups.enumerations")

	fmt.Println(enumConfig)
	
	for _, name := range enumConfig {
		enum := g.GetEnumeration(name)
		enums = append(enums, enum)
	}

	log.Println(enums)

	return enums
}

// GetEnumeration get a single enumeration by name from the API using introspection
func (g *GraphCMS) GetEnumeration(name string) map[string]interface{} {
	var query string = `query EnumerationValues {
		__type(name: "%v") {
		  kind
		  name
		  description
		  enumValues {
			name
			description
		  }
		}
	  }`
	nodeTypes, err := g.CallGraphAPI(fmt.Sprintf(query, name), "{}")

	// Handle any returned errors
	if err != nil {
		log.Printf("Error getting enumerations from api: \n%v", err)
	}

	return nodeTypes.Data
}
