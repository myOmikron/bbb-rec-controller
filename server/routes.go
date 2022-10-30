package server

import (
	"github.com/labstack/echo/v4"

	"github.com/myOmikron/bbb-rec-controller/handler"
	"github.com/myOmikron/bbb-rec-controller/models"
)

func defineRoutes(e *echo.Echo, conf *models.Config) {
	api := handler.Wrapper{
		Config: conf,
	}

	e.GET("/isRecordingRunning", api.IsRecordingRunning)
	e.GET("/stopRecording", api.StopRecording)
}
