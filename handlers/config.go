package handlers

import (
	"net/http"

	"github.com/byuoitav/common/structs"
	"github.com/labstack/echo"
)

// GetConfig pulls the configuration for the device and returns it to the caller
func GetConfig(config structs.LabConfig) func(echo.Context) error {

	return func(ctx echo.Context) error {

		ctx.JSON(http.StatusOK, config)

		return nil
	}
}
