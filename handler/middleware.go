package handler

import "github.com/gin-gonic/gin"

func Middleware(gc *gin.Context) {
	gc.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	gc.Next()
}
