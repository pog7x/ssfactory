package firefoxdrv

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
)

type ServiceFirefox struct {
	FirefoxCaps firefox.Capabilities
	*selenium.Service
}

type InitFirefox struct {
	WebdriverPort     uint16
	FirefoxBinaryPath string
	GeckodriverPath   string
	FirefoxArgs       []string
	Options           []selenium.ServiceOption
}

func NewServiceFirefox(init InitFirefox) (*ServiceFirefox, error) {
	service, err := selenium.NewGeckoDriverService(init.GeckodriverPath, int(init.WebdriverPort), init.Options...)
	if err != nil {
		return nil, err
	}

	firefoxCaps := firefox.Capabilities{
		Binary: init.FirefoxBinaryPath,
		Args:   init.FirefoxArgs,
	}

	return &ServiceFirefox{
		FirefoxCaps: firefoxCaps,
		Service:     service,
	}, nil
}
