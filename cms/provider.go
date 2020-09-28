/*
Copyright © 2020 John Arroyo

CMS Provider interface

All CMS implementations should implement this interface.
*/

package cms

// Provider interface
type Provider interface {
	GetSchema(name string) string
	GetSchemas() string
	GetNodes() []interface{}
	DownloadContent()
}
