package server

import (
	"github.com/labstack/echo/v4"
	"github.com/myOmikron/bbb-rec-controller/modules/bigbluebutton"

	"github.com/myOmikron/bbb-rec-controller/handler"
	"github.com/myOmikron/bbb-rec-controller/models"
)

func defineRoutes(e *echo.Echo, conf *models.Config, bbb *bigbluebutton.BBB) {
	api := handler.Wrapper{
		Config: conf,
		BBB:    bbb,
	}

	e.GET("/isRecordingRunning", api.IsRecordingRunning)
	e.GET("/stopRecording", api.StopRecording)
}
