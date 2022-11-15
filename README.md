# ssfactory (screenshot factory)

[![GoDoc](https://pkg.go.dev/badge/github.com/pog7x/ssfactory)](https://pkg.go.dev/github.com/pog7x/ssfactory)
[![Build Status](https://github.com/pog7x/ssfactory/actions/workflows/go.yml/badge.svg)](https://github.com/pog7x/ssfactory/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/pog7x/ssfactory/blob/master/LICENSE)

## Example of factory initializing

```go
package main

import (
	"github.com/pog7x/ssfactory"
)

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
			WorkersCount: 8,
		},
	)
	if err != nil {
		panic("initialization screenshot factory error")
	}

	defer stopFunc()

	var maximize string
	go f.MakeScreenshot(
		ssfactory.MakeScreenshotPayload{
			URL:            "https://pog7x.github.io/evklid/",
			DOMElementBy:   ssfactory.ByTagName,
			DOMElementName: "body",
			Scroll:         true,
			BytesHandler:   someActionsOnBytes,
			MaximizeWindow: &maximize,
		},
	)
}
```