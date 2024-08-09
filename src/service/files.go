package service

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Transfer serves the content of a file as the HTTP response body.
func Transfer(path, root string, w http.ResponseWriter) error {
	directory, err := filepath.Abs(root)
	if err != nil {
		http.NotFound(w, nil)
		return err
	}

	file, err := filepath.Abs(filepath.Join(directory, path))
	if err != nil || !strings.HasPrefix(file, directory) || !fileExists(file) {
		http.NotFound(w, nil)
		return err
	}

	contentType := "text/plain"
	switch filepath.Ext(file) {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpeg", ".jpg":
		contentType = "image/jpeg"
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, nil, file)
	return nil
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
