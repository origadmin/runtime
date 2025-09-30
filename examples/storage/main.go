package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/origadmin/runtime/interfaces/storage"
)

func main() {
	// Setup base path for storage components
	baseDir, err := os.MkdirTemp("", "storage_server_test")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(baseDir) // Clean up after test

	// Initialize the storage service
	fs, err := NewLocalStorage(baseDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Get the directory of the main.go file
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// Load HTML templates
	r.LoadHTMLGlob(filepath.Join(basepath, "templates/*"))

	// Serve static files (e.g., Bootstrap CSS/JS)
	r.Static("/static", filepath.Join(basepath, "static"))

	// HTTP Handlers
	r.GET("/", indexHandler(fs))
	r.POST("/upload", uploadHandler(fs))
	r.GET("/download", downloadHandler(fs))
	r.POST("/download-zip", downloadZipHandler(fs))
	r.POST("/mkdir", mkdirHandler(fs))
	r.POST("/delete", deleteHandler(fs)) // New handler for delete
	r.POST("/rename", renameHandler(fs)) // New handler for rename

	port := ":8080"
	log.Printf("Starting Gin HTTP server on port %s", port)
	fmt.Println("Application started. Checking logs...")
	log.Fatal(r.Run(port))
}

// indexHandler displays the file list for a given path
func indexHandler(fs storage.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Query("path")
		if currentPath == "" {
			currentPath = "/"
		}
		// Normalize currentPath to use forward slashes for path.Dir
		currentPath = filepath.ToSlash(currentPath)
		currentPath = path.Clean(currentPath) // Normalize path

		files, err := fs.List(currentPath)
		var data TemplateData
		if err != nil {
			data.Error = fmt.Sprintf("Error listing files: %v", err)
		} else {
			parentPath := ""
			if currentPath != "/" {
				parentPath = path.Dir(currentPath)
				if parentPath == "." { // Handle case where path.Dir("/") returns "."
					parentPath = "/"
				}
			}
			// --- DEBUG START ---
			log.Printf("DEBUG: currentPath = %s, calculated parentPath = %s", currentPath, parentPath)
			// --- DEBUG END ---
			data = TemplateData{
				CurrentPath: currentPath,
				ParentPath:  parentPath,
				Files:       files,
				PathParts:   generatePathParts(currentPath),
			}

			// Convert PathParts to JSON for JavaScript
			pathPartsJSON, err := json.Marshal(data.PathParts)
			if err != nil {
				log.Printf("Error marshalling PathParts to JSON: %v", err)
				data.Error = fmt.Sprintf("Error processing path: %v", err)
			} else {
				c.HTML(http.StatusOK, "index.html", gin.H{
					"CurrentPath":   data.CurrentPath,
					"ParentPath":    data.ParentPath,
					"Files":         data.Files,
					"Message":       data.Message,
					"Error":         data.Error,
					"PathParts":     data.PathParts,
					"PathPartsJSON": string(pathPartsJSON),
				})
				return
			}
		}

		c.HTML(http.StatusOK, "index.html", data)
	}
}

// uploadHandler handles file uploads
func uploadHandler(fs storage.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.PostForm("currentPath") // Get current path from form
		if currentPath == "" {
			currentPath = "/"
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.Redirect(http.StatusFound, "/?path="+currentPath)
			return
		}

		fileName := c.PostForm("fileName")
		if fileName == "" {
			c.Redirect(http.StatusFound, "/?path="+currentPath)
			return
		}
		fileName = path.Clean(fileName) // Normalize the full path received from frontend

		src, err := fileHeader.Open()
		if err != nil {
			c.Redirect(http.StatusFound, "/?path="+currentPath)
			return
		}
		defer src.Close()

		// Use fileName directly as it's already the full path from the frontend
		if err := fs.Write(fileName, src, fileHeader.Size); err != nil {
			log.Printf("Failed to save file: %v", err)
			c.Redirect(http.StatusFound, "/?path="+currentPath)
			return
		}

		c.Redirect(http.StatusFound, "/?path="+currentPath)
	}
}

// downloadHandler handles single file downloads
func downloadHandler(fs storage.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := c.Query("path")
		if fileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
			return
		}
		fileName = path.Clean(fileName) // Normalize path

		file, err := fs.Read(fileName)
		if err != nil {
			if os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open file: %v", err)})
			}
			return
		}
		defer file.Close()

		// Set content type and disposition
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(fileName)))

		// Stream the file content
		_, err = io.Copy(c.Writer, file)
		if err != nil {
			log.Printf("Error serving file %s: %v", fileName, err)
		}
	}
}

// downloadZipHandler handles downloading multiple selected files as a ZIP archive
func downloadZipHandler(fs storage.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		selectedFiles := c.PostFormArray("selectedFiles")
		if len(selectedFiles) == 0 {
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.Writer.Header().Set("Content-Type", "application/zip")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

		zipWriter := zip.NewWriter(c.Writer)
		defer zipWriter.Close()

		// Helper function to add a file to the zip; defined outside the loop for clarity.
		addFile := func(filePath string) error {
			file, err := fs.Read(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			info, err := fs.Stat(filePath)
			if err != nil {
				return err
			}

			header := &zip.FileHeader{
				Name:     strings.TrimPrefix(filePath, "/"),
				Method:   zip.Deflate,
				Modified: info.ModTime,
			}

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, file)
			return err
		}

		for _, p := range selectedFiles {
			filePath := path.Clean(p) // Normalize path
			if err := addFile(filePath); err != nil {
				log.Printf("Failed to add file %s to zip: %v", filePath, err)
				continue // Continue to the next file
			}
		}
	}
}

// mkdirHandler handles directory creation
func mkdirHandler(fs storage.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentPath := c.PostForm("parentPath")
		dirName := c.PostForm("dirName")

		if dirName == "" {
			// Handle error: directory name is required
			c.Redirect(http.StatusFound, "/?path="+parentPath)
			return
		}

		newPath := path.Join(parentPath, dirName)
		if err := fs.Mkdir(newPath); err != nil {
			log.Printf("Failed to create directory: %v", err)
			// Handle error, maybe with a flash message if implemented
		}

		c.Redirect(http.StatusFound, "/?path="+parentPath)
	}
}

// deleteHandler handles file/directory deletion
func deleteHandler(fs storage.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetPath := c.PostForm("path") // Assuming path is sent via form for POST
		if targetPath == "" {
			c.Redirect(http.StatusFound, "/") // Redirect to root if no path provided
			return
		}
		targetPath = filepath.ToSlash(targetPath) // Normalize path
		targetPath = path.Clean(targetPath)       // Clean path

		currentPath := path.Dir(targetPath) // Get current directory for redirection
		if currentPath == "." {
			currentPath = "/"
		}

		if err := fs.Delete(targetPath); err != nil {
			log.Printf("Failed to delete %s: %v", targetPath, err)
			// Optionally, add error message to template data
		}
		c.Redirect(http.StatusFound, "/?path="+currentPath)
	}
}

// renameHandler handles file/directory renaming
func renameHandler(fs storage.FileStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		oldPath := c.PostForm("oldPath")
		newFileName := c.PostForm("newPath") // User input: e.g., new_file.txt

		if oldPath == "" || newFileName == "" {
			c.Redirect(http.StatusFound, "/") // Redirect to root if paths are missing
			return
		}

		oldPath = filepath.ToSlash(oldPath) // Normalize path
		oldPath = path.Clean(oldPath)       // Clean path

		// Get the directory of the old file
		oldDir := path.Dir(oldPath)
		// Construct the full new path within the same directory
		newFullPath := path.Join(oldDir, newFileName) // Use newFileName here

		if err := fs.Rename(oldPath, newFullPath); err != nil {
			log.Printf("Failed to rename %s to %s: %v", oldPath, newFullPath, err)
			// Optionally, add error message to template data
		}

		// Redirect back to the directory where the file was (or is now)
		currentPath := oldDir
		if currentPath == "." { // Handle case where path.Dir("/") returns "."
			currentPath = "/"
		}
		c.Redirect(http.StatusFound, "/?path="+currentPath)
	}
}

func generatePathParts(currentPath string) []PathPart {
	var parts []PathPart

	// Normalize currentPath to use forward slashes
	normalizedPath := strings.ReplaceAll(currentPath, "\\", "/")

	// Always add root
	parts = append(parts, PathPart{
		Name:   "ROOT",
		Path:   "/",
		IsLast: normalizedPath == "/",
	})

	if normalizedPath == "/" {
		return parts
	}

	// Split path and build parts
	split := strings.Split(strings.TrimPrefix(normalizedPath, "/"), "/")

	cumulativePath := "/"
	for i, part := range split {
		if part == "" {
			continue
		}

		cleanPart := strings.ReplaceAll(part, "\\", "/")

		if cumulativePath == "/" {
			cumulativePath += cleanPart
		} else {
			cumulativePath += "/" + cleanPart
		}

		newPart := PathPart{
			Name:   cleanPart,
			Path:   cumulativePath,
			IsLast: i == len(split)-1,
		}
		parts = append(parts, newPart)
	}

	return parts
}
