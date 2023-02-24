# ssfactory (screenshot factory)

[![GoDoc](https://pkg.go.dev/badge/github.com/pog7x/ssfactory)](https://pkg.go.dev/github.com/pog7x/ssfactory)
[![Build Status](https://github.com/pog7x/ssfactory/actions/workflows/go.yml/badge.svg)](https://github.com/pog7x/ssfactory/actions/workflows/go.yml)
[![License](https://img.shields.io/github/license/pog7x/ssfactory)](https://github.com/pog7x/ssfactory/blob/master/LICENSE)

## Example of factory initialization and making screenshot

```go
package main

import "github.com/pog7x/ssfactory"

func someActionsOnBytes(screenshotBytes []byte) error {
	// Some actions on bytes (encode to png, send via http/rmq etc.)
	return nil
}

func main() {
	f, stopFunc, err := ssfactory.NewFactory(
		ssfactory.InitFactory{
			WebdriverPort:     8080,
			UseBrowser:        "firefox",
			FirefoxBinaryPath: "/usr/bin/firefox",
			GeckodriverPath:   "./dependencies/geckodriver",
			FirefoxArgs: []string{
				"--headless",
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--width=1920",
			},
		},
	)
	if err != nil {
		panic("initialization screenshot factory error")
	}

	defer stopFunc()

	var maximize string

	f.MakeScreenshot(
		ssfactory.MakeScreenshotPayload{
			URL:            "https://github.com",
			DOMElementBy:   ssfactory.ByTagName,
			DOMElementName: "body",
			Scroll:         true,
			BytesHandler:   someActionsOnBytes,
			MaximizeWindow: &maximize,
		},
	)
}
```