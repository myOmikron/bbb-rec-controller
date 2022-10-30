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
	if !w.BBB.IsValid("isRecordingRunning", &q) {
		return c.XML(400, errorResponse{
			ReturnCode: "FAILED",
			MessageKey: "ChecksumIncorrect",
			Message:    "The provided checksum was incorrect",
		})
	}

	isRunning, err := w.BBB.IsRecordingRunning(w.SeleniumPool, form.MeetingID)
	if err != nil {
		c.Logger().Error(err)
		return c.XML(500, errorResponse{
			ReturnCode: "FAILED",
			MessageKey: "InternalServerError",
			Message:    err.Error(),
		})
	}

	return c.XML(200, isRecordingRunningResponse{
		ReturnCode: "SUCCESS",
		Running:    isRunning,
	})
}
