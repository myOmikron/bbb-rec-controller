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
	Username     string
}

type Selenium struct {
	SeleniumPath     string
	GeckoDriverPath  string
	FirefoxPath      string
	ChromedriverPath string
	ChromiumPath     string
	UseChromium      bool
	DisableHeadless  bool
	InstanceCount    uint
}

type Config struct {
	Server        Server
	BigBlueButton BigBlueButton
	Selenium      Selenium
}
