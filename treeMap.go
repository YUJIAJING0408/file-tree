package fileTree

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

/*
@Date:
@Auth: YUJIAJING
@Desp:
*/

//go:embed static
var staticFiles embed.FS

func NewTreeMapHttp(addr, tmpPath string) {
	r := gin.Default()
	subFS, _ := fs.Sub(staticFiles, "static")
	r.StaticFS("/static", http.FS(subFS))
	r.Static("/tmp", tmpPath)
	r.LoadHTMLFS(http.FS(staticFiles), "static/treeMap.html")

	r.GET("/treeMap", func(c *gin.Context) {
		fileName := c.Query("file_name")
		if fileName == "" {
			c.JSON(http.StatusNotFound, gin.H{})
		}
		c.HTML(http.StatusOK, "treeMap.html", gin.H{
			"version": VERSION,
			"path":    fmt.Sprintf("tmp/%s.json", fileName),
		})
	})
	err := r.Run(addr)
	if err != nil {
		return
	}
}
