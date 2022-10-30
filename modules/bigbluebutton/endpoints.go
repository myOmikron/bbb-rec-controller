package bigbluebutton

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrHTTP      = errors.New("http error, see logs")
	ErrDecodeXML = errors.New("couldn't parse xml")
)

type MeetingInfo struct {
	XMLName    xml.Name `xml:"response"`
	ReturnCode *string  `xml:"returncode"`
	MessageKey *string  `xml:"messageKey"`
	Recording  *bool    `xml:"recording"`
}

// GetMeetingInfo produces a request ready url to GET `getMeetingInfo` with.
func (bbb *BBB) GetMeetingInfo(meetingID string) (*MeetingInfo, error) {
	values := make(url.Values)
	values.Set("meetingID", meetingID)
	bbb.AddChecksum("getMeetingInfo", &values)

	request := *bbb.BaseUrl
	request.RawQuery = values.Encode()
	request.Path = "/bigbluebutton/api/getMeetingInfo"

	response, err := http.Get(request.String())
	if err != nil {
		bbb.Logger.Error(err)
		return nil, ErrHTTP
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		bbb.Logger.Error(err)
		return nil, ErrHTTP
	}

	meetingInfo := &MeetingInfo{}
	err = xml.Unmarshal(body, meetingInfo)
	if err != nil {
		return nil, ErrDecodeXML
	}

	return meetingInfo, nil
}

// Join produces a request ready url to GET `join` with.
func (bbb *BBB) Join(fullName, meetingID string) string {
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
