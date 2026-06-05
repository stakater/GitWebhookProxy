package proxy

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/stakater/GitWebhookProxy/pkg/assets"
)

// serveAssets serves static assets from the embedded filesystem
func (p *Proxy) serveAssets(w http.ResponseWriter, r *http.Request) {
	// Extract the filename from the URL path
	filename := strings.TrimPrefix(r.URL.Path, "/assets/")

	// Prevent directory traversal attacks
	if strings.Contains(filename, "..") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Get the file data
	data, err := assets.GetWebAssetData(filename)
	if err != nil {
		log.Printf("Error reading asset file %s: %v", filename, err)
		http.NotFound(w, r)
		return
	}

	// Set the content type based on file extension
	contentType := "application/octet-stream"
	switch filepath.Ext(filename) {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".html":
		contentType = "text/html"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
