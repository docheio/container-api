package handler

import (
	"context"
	"encoding/json"

	"github.com/docheio/container-api/handshake"
	"github.com/gin-gonic/gin"
)

func (c *Config) Get(gc *gin.Context) {
	var response handshake.Response

	if id := gc.Param("id"); id == "" {
		gc.Status(403)
		return
	}
	
}
