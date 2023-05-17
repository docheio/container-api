package main

import (
	"github.com/docheio/container-api/handler"
	"github.com/docheio/container-api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	var flag utils.Flag
	var router *gin.Engine
	var v1 *gin.RouterGroup

	flag.Init()
	if *flag.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	handlerV1 := handler.Config{
		Uniquekey: flag.Uniquekey,
		Namespace: flag.Namespace,
		Image:     flag.Image,
	}

	handlerV1.Init()
	router = gin.Default()
	v1 = router.Group("/v1")

	v1.Use(handler.Middleware)
	v1.GET("/", handlerV1.Get)

}
