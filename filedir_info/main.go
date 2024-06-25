package main

import (
	"filedir_info/file"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Item represents a file system item with its type and path
type Item struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func main() {
	r := gin.Default()

	// Define a route for getting directories and disks
	r.GET("/", func(c *gin.Context) {
		// Get the path from query parameters
		path := c.Query("path")

		// Get directories and disks in the specified path
		items, err := file.NewFileInfo(file.FileOption{
			Path:       path,
			Search:     "",
			ContainSub: false,
			Expand:     true,
			Dir:        false,
			ShowHidden: true,
			Page:       1,
			PageSize:   100,
			SortBy:     "name",
			SortOrder:  "ascending",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	// Run the Gin server
	r.Run(":8080")
}
