package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (config *Config) checkInternalServerError(gc *gin.Context, err error) {
	if err != nil {
		gc.Writer.WriteHeader(http.StatusInternalServerError)
		gc.Writer.WriteString("Internal Server Error")
	}
}
