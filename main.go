/* ********************************************************************************************************** */
/*                                                                                                            */
/*                                                     :::::::::  ::::::::   ::::::::   :::    ::: :::::::::: */
/* main.go                                            :+:    :+: :+:    :+: :+:    :+: :+:    :+: :+:         */
/*                                                   +:+    +:+ +:+    +:+ +:+        +:+    +:+ +:+          */
/* By: ES-Yukun <yukun@doche.io>                    +#+    +:+ +#+    +:+ +#+        +#++:++#++ +#++:++#      */
/*                                                 +#+    +#+ +#+    +#+ +#+        +#+    +#+ +#+            */
/* Created: 2023/08/09 23:48:06 by ES-Yukun       #+#    #+# #+#    #+# #+#    #+# #+#    #+# #+#             */
/* Updated: 2023/08/09 23:48:07 by ES-Yukun      #########  ########   ########  ###    ### ##########.io.    */
/*                                                                                                            */
/* ********************************************************************************************************** */

package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/docheio/container-api/handler"
	"github.com/docheio/container-api/utils"
)

func main() {
	flags := utils.Flags{}
	flags.Init()

	if *flags.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	handler := handler.Handler{}
	handler.Init()
	handler.Namespace = *flags.Namespace
	handler.Uniquekey = *flags.Uniquekey
	handler.Image = *flags.Image

	router := gin.Default()

	routerV1 := router.Group("/v1")
	routerV1.POST("/", handler.Create)
	routerV1.GET("/", handler.GetAll)
	routerV1.GET("/:id", handler.GetOne)
	routerV1.PUT("/:id", handler.Update)
	routerV1.DELETE("/:id", handler.Delete)

	log.Println("Server Start")
	if err := router.Run(*flags.Address); err != nil {
		log.Fatalln(err)
	}
}
