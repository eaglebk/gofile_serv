package files

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"space.eagle2000/fileserv/config"
)

var (
	dir, _   = os.Getwd()
	conf, _  = config.NewConfig(path.Join(dir, "config.yml"))
	fileInfo *os.FileInfo
	err      error
)

// GetFile - returns the file
func GetFile(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "" {
		c.String(http.StatusBadRequest, "Get 'filename' not specified in url.")
		return
	}

	pathOfFile := filepath.Join(conf.Server.FilePath, filename)

	openfile, err := os.Open(pathOfFile)
	defer openfile.Close()

	if err != nil {
		c.String(http.StatusBadRequest, "File not found.")
		return
	}

	fileHeader := make([]byte, 512)
	openfile.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	fileStat, _ := openfile.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", fileContentType)
	c.Header("Content-Length", fileSize)

	openfile.Seek(0, 0)
	io.Copy(c.Writer, openfile)
	return
}

// UploadFile - upload the file
func UploadFile(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	filename := header.Filename

	os.MkdirAll(conf.Server.FilePath, os.ModePerm)
	out, err := os.Create(filepath.Join(conf.Server.FilePath, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	c.String(http.StatusOK, "File uploaded")
}

// UpdateFile - replace the file
func UpdateFile(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	pathOfFile := filepath.Join(conf.Server.FilePath, filename)

	// detect if file exists
	var _, errP = os.Stat(pathOfFile)

	if os.IsNotExist(errP) {
		c.String(http.StatusBadRequest, "File not exists")
		return
	}

	os.MkdirAll(conf.Server.FilePath, os.ModePerm)
	out, err := os.Create(filepath.Join(conf.Server.FilePath, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	c.String(http.StatusOK, "File updated")
}

// DeleteFile - delete the file
func DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	fmt.Println(filename)

	var err = os.Remove(filepath.Join(conf.Server.FilePath, filename))
	if err != nil {
		c.String(http.StatusBadRequest, "File not found")
		return
	}

	c.String(http.StatusOK, "File was deleted")

}
