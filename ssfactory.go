package ssfactory

import (
	"fmt"
	"time"

	"github.com/pog7x/ssfactory/internal/chromedrv"
	"github.com/pog7x/ssfactory/internal/firefoxdrv"
	"github.com/pog7x/ssfactory/internal/workerpool"

	"github.com/tebeka/selenium"
)

const (
	FirefoxName = "firefox"
	ChromeName  = "chrome"
)

const (
	ByID              = "id"
	ByXPATH           = "xpath"
	ByLinkText        = "link text"
	ByPartialLinkText = "partial link text"
	ByName            = "name"
	ByTagName         = "tag name"
	ByClassName       = "class name"
	ByCSSSelector     = "css selector"
)

type Factory struct {
	wp *workerpool.WorkerPool

	firefoxService *firefoxdrv.ServiceFirefox
	chromeService  *chromedrv.ServiceChrome

	capabilities selenium.Capabilities

	webdriverPort uint16
	webdriver     *selenium.WebDriver

	urlBase string
}

type InitFactory struct {
	WebdriverPort uint16

	UseBrowser string

	FirefoxBinaryPath string
	FirefoxArgs       []string
	GeckodriverPath   string

	ChromeBinaryPath string
	ChromeArgs       []string
	ChromedriverPath string

	WorkersCount uint8
}

func NewFactory(init InitFactory) (*Factory, func(), error) {
	f := &Factory{
		capabilities:  selenium.Capabilities{"browserName": init.UseBrowser},
		webdriverPort: init.WebdriverPort,
		wp:            workerpool.NewWP(init.WorkersCount),
	}

	switch init.UseBrowser {
	case FirefoxName:
		firefoxSvc, err := firefoxdrv.NewServiceFirefox(
			firefoxdrv.InitFirefox{
				WebdriverPort:     init.WebdriverPort,
				FirefoxBinaryPath: init.FirefoxBinaryPath,
				GeckodriverPath:   init.GeckodriverPath,
				FirefoxArgs:       init.FirefoxArgs,
			},
		)
		if err != nil {
			return nil, func() {}, err
		}
		f.firefoxService = firefoxSvc
		f.capabilities.AddFirefox(f.firefoxService.FirefoxCaps)
	case ChromeName:
		chromeSvc, err := chromedrv.NewServiceChrome(
			chromedrv.InitChrome{
				WebdriverPort:    init.WebdriverPort,
				ChromeBinaryPath: init.ChromeBinaryPath,
				ChromedriverPath: init.ChromedriverPath,
				ChromeArgs:       init.ChromeArgs,
			},
		)
		if err != nil {
			return nil, func() {}, err
		}
		f.chromeService = chromeSvc
		f.capabilities.AddChrome(f.chromeService.ChromeCaps)
		f.urlBase = "/wd/hub"
	default:
		return nil, func() {}, fmt.Errorf("specified invalid browser name (%s)", init.UseBrowser)
	}

	if err := f.runFactory(); err != nil {
		return nil, func() {}, err
	}

	return f, f.stopFactory, nil
}

func (f *Factory) runFactory() error {
	f.wp.Run()
	wd, err := selenium.NewRemote(f.capabilities, fmt.Sprintf("http://localhost:%d%s", f.webdriverPort, f.urlBase))
	if err != nil {
		return err
	}
	f.webdriver = &wd
	return nil
}

func (f *Factory) stopFactory() {
	if f.webdriver != nil {
		if err := (*f.webdriver).Quit(); err != nil {

		}
	}
	if f.firefoxService != nil {
		if err := f.firefoxService.Stop(); err != nil {

		}
	}
	if f.chromeService != nil {
		if err := f.chromeService.Stop(); err != nil {
		}
	}
	if f.wp != nil {
		f.wp.Stop()
	}
}

type bytesHandler func([]byte) error

type MakeScreenshotPayload struct {
	URL            string
	DOMElementBy   string
	DOMElementName string
	Scroll         bool
	MaximizeWindow *string
	Timeout        time.Duration
	BytesHandler   bytesHandler
}

func (f *Factory) MakeScreenshot(p MakeScreenshotPayload) {
	f.wp.Do(func() {
		if err := (*f.webdriver).Get(p.URL); err != nil {
			return
		}

		if p.MaximizeWindow != nil {
			if err := (*f.webdriver).MaximizeWindow(*p.MaximizeWindow); err != nil {
				return
			}
		}

		if p.Timeout != 0 {
			if err := (*f.webdriver).WaitWithTimeout(
				func(wd selenium.WebDriver) (bool, error) { return false, nil },
				p.Timeout,
			); err != nil {
				return
			}
		}

		elem, err := (*f.webdriver).FindElement(p.DOMElementBy, p.DOMElementName)
		if err != nil {
			return
		}

		screenshotBytes, err := elem.Screenshot(p.Scroll)
		if err != nil {
			return
		}

		p.BytesHandler(screenshotBytes)
	})
}
