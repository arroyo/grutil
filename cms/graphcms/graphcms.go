/*
Copyright © 2020 John Arroyo

cms graphcms package

Get and persist both schemas and content from GraphCMS
*/

package graphcms

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// GraphCMS struct
type GraphCMS struct {
	url          interface{}
	key          interface{}
	path         string
	folder       string
	structure    map[string]string
	stage        string
	NodeTypes    []string
	SpecialCases map[string]string
	Debug        bool
	Exceptions   map[string]string
	RenderData   GraphResponse
}

// AssetNode simple content struct
type AssetNode struct {
	TypeName string `json:"__typename"`
	ID       string `json:"id"`
	Handle   string `json:"handle"`
	URL      string `json:"url"`
}

// Init initialize config
func (g *GraphCMS) Init(url interface{}, key interface{}, stage interface{}, path interface{}) {
	g.url = url
	g.key = key
	g.path = fmt.Sprintf("%v", path)
	g.stage = fmt.Sprintf("%v", stage)
	g.SpecialCases = map[string]string{
		"RichText": "%v { raw html markdown text }\n",
		"Asset":    "%v { id }\n",
	}
	g.Debug = false
	g.Exceptions = viper.GetStringMapString("exceptions.plurals")
	if g.Debug {
		log.Println(g.Exceptions)
	}
}

// Set the debug flag
func (g *GraphCMS) SetDebug(debug bool) {
	g.Debug = debug
}

// GetNodes from the cms
func (g *GraphCMS) GetNodes() map[string]interface{} {
	// Get Node Types, loop through each type and then pull all content
	g.NodeTypes = g.GetNodeTypes()
	allNodes := make(map[string]interface{})

	for index := range g.NodeTypes {
		log.Printf("Get %v nodes", g.NodeTypes[index])
		allNodes[g.NodeTypes[index]] = g.GetAllNodesByType(g.NodeTypes[index])
	}

	return allNodes
}

// subfieldException check for fields that cannot be accounted for by introspection
func (g *GraphCMS) subfieldException(name string) string {
	format := ""

	// rgba is a normalized field with no distinguisging introspection data
	// distance is based on a map input field that is irrelavant for a backup
	switch name {
	case "rgba":
		format = "rgba { r g b a }\n"
	case "distance":
		format = ""
	default:
		format = name + "\n"
	}

	return format
}

// IsNodeType check to see if name is listed as one of the node types
func (g *GraphCMS) IsNodeType(name string) bool {
	for _, nodeType := range g.NodeTypes {
		if strings.ToLower(name) == strings.ToLower(nodeType) {
			return true
		}
	}

	return false
}

// IsSpecialCase field type
func (g *GraphCMS) IsSpecialCase(name string) bool {
	for key := range g.SpecialCases {
		if key == name {
			return true
		}
	}

	return false
}

// subfieldFormat will apply special rules to handle fields that both
// can and cannot be accounted for based on introspection.
func (g *GraphCMS) subfieldFormat(field NodeSubfield) string {
	// To avoid recursive nesting, check for one of the known schema types and simply add an id if a match
	if g.IsNodeType(field.Name) {
		return field.Name + " { id }\n"
	}

	// If fields are present pull them and return subquery
	if len(field.Type.Fields) > 0 {
		subfields := ""
		for _, f := range field.Type.Fields {
			if len(f.Args) > 1 {
				subfields += f.Name + " { id } \n"
			} else {
				subfields += g.subfieldException(f.Name)
			}
		}
		return fmt.Sprintf("%v { %v }\n", field.Name, subfields)
	}

	// If there are multiple args it indicates we need a refrence id
	if len(field.Args) > 1 {
		return field.Name + " { id }\n"
	}

	//  Else apply subfield exceptions and return
	return g.subfieldException(field.Name)
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

	// Handle any returned errors
	if err != nil {
		log.Fatalf("error processing node fields: \n%v", err)
	}

	// Query builder, create GraphQL query to pull all nodes and content fields
	query = `query {
		%v {
			__typename
			%v
		}
	}`
	for index := range nodeFields.Type.Fields {
		// Is it one of our custom node types
		if g.IsNodeType(nodeFields.Type.Fields[index].Type.Name) {
			fieldsQuery += nodeFields.Type.Fields[index].Name + " { id }\n"

			// Is it a union field?
		} else if nodeFields.Type.Fields[index].Type.Kind == "UNION" {
			fields := ""
			for i := range nodeFields.Type.Fields[index].Type.PossibleTypes {
				fields += fmt.Sprintf("... on %v { id }\n", nodeFields.Type.Fields[index].Type.PossibleTypes[i].Name)
			}
			fieldsQuery += fmt.Sprintf("%v { %v } \n", nodeFields.Type.Fields[index].Name, fields)

			// Do we have fields defined?
		} else if len(nodeFields.Type.Fields[index].Type.Fields) > 0 {
			fields := ""
			for _, field := range nodeFields.Type.Fields[index].Type.Fields {
				fields += g.subfieldFormat(field)
			}
			fieldsQuery += fmt.Sprintf("%v { %v } \n", nodeFields.Type.Fields[index].Name, fields)

			// Is it in our special case map?
		} else if g.IsSpecialCase(nodeFields.Type.Fields[index].Type.OfType.Name) {
			fieldsQuery += fmt.Sprintf(g.SpecialCases[nodeFields.Type.Fields[index].Type.OfType.Name], nodeFields.Type.Fields[index].Name)

			// Do we still need this?
		} else if len(nodeFields.Type.Fields[index].Args) > 1 {
			fieldsQuery += fmt.Sprintf("%v { id } \n", nodeFields.Type.Fields[index].Name)

			// Default
		} else {
			fieldsQuery += nodeFields.Type.Fields[index].Name + "\n"
		}
	}

	nodes := g.Pluralize(name)
	query = fmt.Sprintf(query, nodes, fieldsQuery)
	allNodes, err := g.CallGraphAPI(query, "{}")

	if err != nil {
		log.Fatalf("error pulling all %s nodes from API: %v", name, err)
	}

	// Normalize name structure
	allNodes.Data["nodes"] = allNodes.Data[nodes]
	delete(allNodes.Data, nodes)

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
		log.Fatal("error getting Node types from GraphCMS API: \n%v", err)
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
		log.Fatal("error parsing node types from response: \n%v", err)
	}

	var allTypes []string
	for index := range nodeTypesResp.Type.PossibleTypes {
		allTypes = append(allTypes, nodeTypesResp.Type.PossibleTypes[index].Name)
	}

	return allTypes
}

// NodeSubfield structure
type NodeSubfield struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Args        []struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		DefaultValue string `json:"defaultValue"`
	} `json:"args"`
	Type struct {
		Fields []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Args        []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"args"`
		} `json:"fields"`
	} `json:"type"`
}

// NodeFields structure
type NodeFields struct {
	Type struct {
		Name   string `json:"name"`
		Kind   string `json:"kind"`
		Fields []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Args        []struct {
				Name         string `json:"name"`
				Description  string `json:"description"`
				DefaultValue string `json:"defaultValue"`
			} `json:"args"`
			Type struct {
				Name          string         `json:"name"`
				Kind          string         `json:"kind"`
				Description   string         `json:"description"`
				Fields        []NodeSubfield `json:"fields"`
				PossibleTypes []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"possibleTypes"`
				OfType struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"ofType"`
				Interfaces []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Typename    string `json:"__typename"`
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
				type {
					fields {
						name
						description
						args {
							name
							description
					  	}
					}
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
		log.Fatal("error getting Node fields from GraphCMS API: \n%v", err)
	}

	return nodeFields.Data
}

// GetSchema returns a single schema as json
func (g *GraphCMS) GetSchema(name string) (interface{}, error) {
	query, queryVars := g.GetSchemaQuery(name)
	response, err := g.CallGraphAPI(query, queryVars)
	if err != nil {
		fmt.Printf("API Failre: %v", err)
	}

	return response.Data["__type"], err
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
						type {
							fields {
								name
								description
							}
						}
					}
				}
			}
		}
	}`
	var queryVars = fmt.Sprintf(`{"nodetype":"%s"}`, name)

	return query, queryVars
}

// GetSchemas retrives the schema and returns json
func (g *GraphCMS) GetSchemas() map[string]interface{} {
	g.NodeTypes = g.GetNodeTypes()
	schemas := make(map[string]interface{})
	var nodeType string
	var err error

	fmt.Printf("NodeTypes: %v\n", g.NodeTypes)

	for index := range g.NodeTypes {
		nodeType = g.NodeTypes[index]
		log.Printf("Get %v schema json", nodeType)
		schemas[nodeType], err = g.GetSchema(nodeType)
		if err != nil {
			log.Printf("error getting %v schema", nodeType)
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
		log.Printf("error getting enumerations from api: \n%v", err)
	}

	type SchemaTypes struct {
		Schema struct {
			Types []struct {
				Name        string `json:"name"`
				Kind        string `json:"kind"`
				Description string `json:"description"`
				EnumValues  []struct {
					Name        string `json:"name"`
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
		log.Fatal("error parsing schema types from response: \n%v", err)
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
func (g *GraphCMS) GetEnumerations() map[string]interface{} {
	enums := make(map[string]interface{})
	enumConfig := viper.GetStringSlice("backups.enumerations")

	for _, name := range enumConfig {
		enums[name] = g.GetEnumeration(name)
	}

	return enums
}

// GetEnumeration get a single enumeration by name from the API using introspection
func (g *GraphCMS) GetEnumeration(name string) interface{} {
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
		log.Printf("error getting enumerations from api: \n%v", err)
	}

	return nodeTypes.Data["__type"]
}
