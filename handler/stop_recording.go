package handler

import (
	"github.com/labstack/echo/v4"
)

type stopRecordingRequest struct {
	MeetingID string `query:"meetingID"`
	Checksum  string `query:"checksum"`
}

type stopRecordingResponse struct {
	ReturnCode    string   `xml:"returncode"`
	StopRecording bool     `xml:"stopRecording"`
	XMLName       struct{} `xml:"response"`
}

func (w *Wrapper) StopRecording(c echo.Context) error {
	form := stopRecordingRequest{}

	if err := c.Bind(&form); err != nil {
		c.Logger().Info(err)
		return c.XML(400, errorResponse{
			ReturnCode: "FAILED",
			MessageKey: "Bad request",
			Message:    err.Error(),
		})
	}

	if form.MeetingID == "" {
		return c.XML(400, errorResponse{
			ReturnCode: "FAILED",
			MessageKey: "MeetingIDMissing",
			Message:    "The parameter meetingID is missing",
		})
	}

	q := c.Request().URL.Query()
	if !w.BBB.IsValid("stopRecording", &q) {
		return c.XML(400, errorResponse{
			ReturnCode: "FAILED",
			Message:    "ChecksumIncorrect",
			MessageKey: "The provided checksum was incorrect",
		})
	}

	running := false

	return c.XML(200, stopRecordingResponse{
		ReturnCode:    "SUCCESS",
		StopRecording: running,
	})
}
