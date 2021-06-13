/*
Copyright © 2021 John Arroyo

CMS Provider interface

All CMS implementations should implement this interface.
*/

package cms

// Provider interface
type Provider interface {
	GetSchema(name string) (interface{}, error)
	GetSchemas() map[string]interface{}
	GetNodes() map[string]interface{}
	DownloadContent()
	DownloadSchemas() error
}
