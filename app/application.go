package app

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"space.eagle2000/fileserv/config"
)

var (
	router = gin.Default()
)

// StartApplication ...
func StartApplication() {
	dir, _ := os.Getwd()
	conf, _ := config.NewConfig(path.Join(dir, "config.yml"))
	mapUrls()
	router.Run(":" + conf.Server.Port)
}
