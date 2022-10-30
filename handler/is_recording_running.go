package handler

import "github.com/labstack/echo/v4"

type isRecordingRunningRequest struct {
	MeetingID string `query:"meetingID"`
	Checksum  string `query:"checksum"`
}

type isRecordingRunningResponse struct {
	ReturnCode string   `xml:"returncode"`
	Running    bool     `xml:"running"`
	XMLName    struct{} `xml:"response"`
}

func (w *Wrapper) IsRecordingRunning(c echo.Context) error {
	form := isRecordingRunningRequest{}

	if err := c.Bind(&form); err != nil {
		c.Logger().Info(err)
		return c.XML(400, errorResponse{
			ReturnCode: "FAILED",
			MessageKey: "BadRequest",
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

	return c.XML(200, isRecordingRunningResponse{
		ReturnCode: "SUCCESS",
	})
}
