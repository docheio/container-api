package main

import (
	"github.com/docheio/container-api/handler"
	"github.com/docheio/container-api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	flags := utils.Flags{}
	flags.Init()

	if *flags.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

}
