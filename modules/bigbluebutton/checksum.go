package bigbluebutton

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"strings"
)

func (bbb *BBB) computeChecksum(endpoint string, values *url.Values) string {
	request := fmt.Sprintf("%s%s%s", endpoint, values.Encode(), bbb.Config.SharedSecret)
	return fmt.Sprintf("%x", sha1.Sum([]byte(request)))
}

// IsValid checks if an url.Values contains a valid checksum for a given endpoint
func (bbb *BBB) IsValid(endpoint string, values *url.Values) bool {
	checksum := values.Get("checksum")
	checksum = strings.ToLower(checksum)
	values.Del("checksum")
	return checksum == bbb.computeChecksum(endpoint, values)
}

// AddChecksum adds the checksum for a specific endpoint to an url.Values
func (bbb *BBB) AddChecksum(endpoint string, values *url.Values) {
	values.Add("checksum", bbb.computeChecksum(endpoint, values))
}
