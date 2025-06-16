package handlers

import (
	"net/http"

	"github.com/byuoitav/common/structs"
	"github.com/gin-gonic/gin"
)

// GetConfig pulls the configuration for the device and returns it to the caller
func GetConfig(config structs.LabConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, config)
	}
}
