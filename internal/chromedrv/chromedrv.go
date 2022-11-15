package chromedrv

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ServiceChrome struct {
	ChromeCaps chrome.Capabilities
	*selenium.Service
}

type InitChrome struct {
	WebdriverPort    uint16
	ChromeBinaryPath string
	ChromedriverPath string
	ChromeArgs       []string
	Options          []selenium.ServiceOption
}

func NewServiceChrome(init InitChrome) (*ServiceChrome, error) {
	service, err := selenium.NewChromeDriverService(init.ChromedriverPath, int(init.WebdriverPort), init.Options...)
	if err != nil {
		return nil, err
	}

	chromeCaps := chrome.Capabilities{
		Path: init.ChromeBinaryPath,
		Args: init.ChromeArgs,
	}

	return &ServiceChrome{
		ChromeCaps: chromeCaps,
		Service:    service,
	}, nil
}
