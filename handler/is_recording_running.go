package handler

import "github.com/labstack/echo/v4"

type isRecordingRunningRequest struct {
	MeetingID string `query:"meetingID"`
	Checksum  string `query:"checksum"`
}

type isRecordingRunningResponse struct {
}

func (w *Wrapper) IsRecordingRunning(c echo.Context) error {

	return c.XML(200, isRecordingRunningResponse{})
}
