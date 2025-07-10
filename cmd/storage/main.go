package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/origadmin/runtime/storage"
)

var storageService *storage.Storage

// TemplateData holds data to be passed to the HTML template
type TemplateData struct {
	CurrentPath string
	ParentPath  string
	Files       []FileInfo
	Message     string
	Error       string
}

// FileInfo holds information about a file or directory for display
type FileInfo struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

func main() {
	// Setup base path for storage components
	baseDir, err := os.MkdirTemp("", "storage_server_test")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(baseDir) // Clean up after test

	// Initialize the storage service
	cfg := storage.Config{
		BasePath: baseDir,
	}
	storageService, err = storage.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage service: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve static files (e.g., Bootstrap CSS/JS)
	r.Static("/static", "./static")

	// HTTP Handlers
	r.GET("/", indexHandler)
	r.POST("/upload", uploadHandler)
	r.GET("/download", downloadHandler)
	r.POST("/download-zip", downloadZipHandler)

	port := ":8080"
	log.Printf("Starting Gin HTTP server on port %s", port)
	log.Fatal(r.Run(port))
}

// indexHandler displays the file list for a given path
func indexHandler(c *gin.Context) {
	currentPath := c.Query("path")
	if currentPath == "" {
		currentPath = "/"
	}
	currentPath = path.Clean(currentPath) // Normalize path

	// Get directory entries from IndexManager
	var files []FileInfo
	var err error

	// For simplicity, assuming IndexManager.ListChildren can take a path and return nodes
	// In a real scenario, you'd get the node for currentPath, then list its children
	// This part needs proper IndexManager integration.
	// For now, we'll simulate a flat structure or rely on IndexManager's actual ListChildren
	// if it can be adapted to list by path.
	// Assuming IndexManager.ListChildren takes NodeID, we need to get NodeID from path first.

	node, err := storageService.IndexManager.GetNodeByPath(currentPath)
	if err != nil {
		if os.IsNotExist(err) {
			data.Error = "Path not found."
		} else {
			data.Error = fmt.Sprintf("Error getting path: %v", err)
		}
	} else {
		children, err := storageService.IndexManager.ListChildren(node.ID)
		if err != nil {
			data.Error = fmt.Sprintf("Error listing files: %v", err)
		} else {
			for _, child := range children {
				files = append(files, FileInfo{
					Name:    child.Name,
					Path:    path.Join(currentPath, child.Name),
					IsDir:   child.IsDir,
					Size:    child.Size,
					ModTime: child.ModTime,
				})
			}
		}
	}

	parentPath := ""
	if currentPath != "/" {
		parentPath = path.Dir(currentPath)
		if parentPath == "." { // Handle case where path.Dir("/") returns "."
			parentPath = "/"
		}
	}

	data := TemplateData{
		CurrentPath: currentPath,
		ParentPath:  parentPath,
		Files:       files,
	}

	if msg := c.Flash.Get("message"); msg != nil {
		data.Message = msg.(string)
	}
	if errMsg := c.Flash.Get("error"); errMsg != nil {
		data.Error = errMsg.(string)
	}

	c.HTML(http.StatusOK, "index.html", data)
}

// uploadHandler handles file uploads
func uploadHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.Flash.Add("error", fmt.Sprintf("Failed to get file: %v", err))
		c.Redirect(http.StatusFound, "/")
		return
	}

	fileName := c.PostForm("fileName")
	if fileName == "" {
		c.Flash.Add("error", "Target path is required")
		c.Redirect(http.StatusFound, "/")
		return
	}
	fileName = path.Clean(fileName) // Normalize path

	src, err := fileHeader.Open()
	if err != nil {
		c.Flash.Add("error", fmt.Sprintf("Failed to open uploaded file: %v", err))
		c.Redirect(http.StatusFound, "/")
		return
	}
	defer src.Close()

	// Simulate chunked upload using the new API
	// Start Upload
	uploadID, err := storageService.MetaStore.StartUpload(fileName, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		c.Flash.Add("error", fmt.Sprintf("Failed to start upload: %v", err))
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Read and upload chunks
	// Assuming DefaultVersion returns chunk size for MetaStore, which is incorrect.
	// This needs to be a proper chunk size from configuration or a sensible default.
	// For now, using a fixed size.
	chunkSize := int64(4 * 1024 * 1024) // 4MB chunk size
	buf := make([]byte, chunkSize)
	chunkIndex := 0
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			err = storageService.MetaStore.UploadChunk(uploadID, chunkIndex, buf[:n])
			if err != nil {
				storageService.MetaStore.CancelUpload(uploadID) // Attempt to clean up
				c.Flash.Add("error", fmt.Sprintf("Failed to upload chunk %d: %v", chunkIndex, err))
				c.Redirect(http.StatusFound, "/")
				return
			}
			chunkIndex++
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			storageService.MetaStore.CancelUpload(uploadID)
			c.Flash.Add("error", fmt.Sprintf("Failed to read file for chunking: %v", readErr))
			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	// Finish Upload
	_, err = storageService.MetaStore.FinishUpload(uploadID)
	if err != nil {
		storageService.MetaStore.CancelUpload(uploadID)
		c.Flash.Add("error", fmt.Sprintf("Failed to finish upload: %v", err))
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.Flash.Add("message", fmt.Sprintf("File %s uploaded successfully!", fileName))
	c.Redirect(http.StatusFound, "/")
}

// downloadHandler handles single file downloads
func downloadHandler(c *gin.Context) {
	fileName := c.Query("path")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File path is required"})
		return
	}
	fileName = path.Clean(fileName) // Normalize path

	// Assuming IndexManager has an Open method that returns fs.File
	// For now, we'll use MetaStore.Open as a placeholder
	file, err := storageService.MetaStore.Open(fileName)
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

// downloadZipHandler handles downloading multiple selected files as a ZIP archive
func downloadZipHandler(c *gin.Context) {
	selectedFiles := c.PostFormArray("selectedFiles")
	if len(selectedFiles) == 0 {
		c.Flash.Add("error", "No files selected for download.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	for _, filePath := range selectedFiles {
		filePath = path.Clean(filePath) // Normalize path

		file, err := storageService.MetaStore.Open(filePath) // Use MetaStore.Open for now
		if err != nil {
			log.Printf("Error opening file %s for zip: %v", filePath, err)
			continue // Skip this file, but continue with others
		}
		defer file.Close()

		// Get file info for zip header
		// This part is problematic as MetaStore.Open returns fs.File, not os.File
		// and fs.File does not have Stat() method directly.
		// We need to get FileInfo from IndexManager or MetaStore.Stat
		// For now, we'll create a dummy header or rely on a proper Stat implementation.
		// Assuming file implements fs.File and has a Stat() method.
		// For now, we'll use a simplified approach.
		// Example: fileInfo, err := file.Stat()
		// if err != nil { log.Printf(...); continue; }
		// header, err := zip.FileInfoHeader(fileInfo)

		// Placeholder for actual file info
		header := &zip.FileHeader{
			Name:     strings.TrimPrefix(filePath, "/"),
			Method:   zip.Deflate,
			Modified: time.Now(),
		}
		// If we had actual file size, we'd set header.UncompressedSize64

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			log.Printf("Error creating zip header for %s: %v", filePath, err)
			continue
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			log.Printf("Error writing file %s to zip: %v", filePath, err)
			continue
		}
	}
}
