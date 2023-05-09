package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (config *Config) checkInternalServerError(gc *gin.Context, err error) bool {
	if err != nil {
		gc.Writer.WriteHeader(http.StatusInternalServerError)
		gc.Writer.WriteString(err.Error())
		return true
	}
	return false
}
