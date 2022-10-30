package models

type Server struct {
	ListenAddress           string
	ListenPort              uint16
	PublicURI               string
	AllowedHosts            []string
	UseForwardedProtoHeader bool
}

type BigBlueButton struct {
	ServerURI    string
	SharedSecret string
}

type Selenium struct {
	SeleniumPath string
	ChromiumPath string
}

type Config struct {
	Server        Server
	BigBlueButton BigBlueButton
	Selenium      Selenium
}
