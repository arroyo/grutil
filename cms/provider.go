/*
Copyright © 2020 John Arroyo

CMS Provider interface

All CMS implementations should implement this interface.
*/

package cms

type Provider interface {
	GetSchema() string
	GetSchemas() string
	GetContent() string
	DownloadContent()
}
