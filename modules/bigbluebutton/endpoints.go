package bigbluebutton

import (
	"github.com/labstack/echo/v4"
	"net/url"
)

// GetMeetingInfo produces a request ready url to GET `getMeetingInfo` with.
func (bbb *BBB) GetMeetingInfo(_ctx echo.Context, meetingID string) string {
	values := make(url.Values)
	values.Set("meetingID", meetingID)
	bbb.AddChecksum("getMeetingInfo", &values)

	request, _ := url.Parse(bbb.Config.ServerURI)
	request.RawQuery = values.Encode()
	request.Path = "/bigbluebutton/api/getMeetingInfo"
	return request.String()
}

// Join produces a request ready url to GET `join` with.
func (bbb *BBB) Join(_ctx echo.Context, fullName, meetingID string) string {
	values := make(url.Values)
	values.Set("fullName", fullName)
	values.Set("meetingID", meetingID)
	values.Set("role", "MODERATOR")
	bbb.AddChecksum("join", &values)

	request, _ := url.Parse(bbb.Config.ServerURI)
	request.RawQuery = values.Encode()
	request.Path = "/bigbluebutton/api/join"
	return request.String()
}
