package main

import (
	"github.com/docheio/container-api/handler"
	"github.com/docheio/container-api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	var flags utils.Flags
	var router *gin.Engine
	var v1 *gin.RouterGroup

	flags = utils.Flags{}
	flags.Init()

	if *flags.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	handler := handler.Config{
		Uniquekey: flags.Uniquekey,
		Namespace: flags.Namespace,
		Image:     flags.Image,
	}

	handler.Init()
	router = gin.Default()
	v1 = router.Group("/v1")

	v1.GET("/", handler.GetAll)
	v1.GET("/:id", handler.GetOne)
	v1.POST("/", handler.Create)
	router.Run(*flags.Address)
}
