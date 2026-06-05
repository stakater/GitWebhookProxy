// Package assets provides access to embedded static assets
package assets

import (
	"embed"
	"io/fs"
	"log"
)

// The embed directive must use patterns relative to the package directory
//
//go:embed web
var webAssets embed.FS

// GetWebAssets returns a filesystem containing the web assets
func GetWebAssets() fs.FS {
	webFS, err := fs.Sub(webAssets, "web")
	if err != nil {
		log.Printf("Error creating web assets sub-filesystem: %v", err)
		return nil
	}
	return webFS
}

// GetWebAssetData returns the content of a web asset file
func GetWebAssetData(filename string) ([]byte, error) {
	return webAssets.ReadFile("web/" + filename)
}
