package main

import "github.com/origadmin/runtime/interfaces/storage"

// TemplateData holds data to be passed to the HTML template
type TemplateData struct {
	CurrentPath   string
	ParentPath    string
	Files         []storage.FileInfo
	Message       string
	Error         string
	PathParts     []PathPart
	PathPartsJSON string // Added for JavaScript consumption
}

// PathPart represents a segment of the current path for breadcrumbs.
type PathPart struct {
	Name   string
	Path   string
	IsLast bool
}
