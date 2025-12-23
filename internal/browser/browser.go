package browser

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices" // Added this import
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
)

func NewBrowser(headless bool) *rod.Browser {
	l := launcher.New().
		Headless(headless).
		Devtools(false).
		Leakless(false).
		Set("start-maximized", ""). // Force Window Maximize
		Set("disable-blink-features", "AutomationControlled")

	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	return browser
}

func CreatePage(browser *rod.Browser) *rod.Page {
	page := stealth.MustPage(browser)
	page.MustEmulate(devices.Clear)
	page.MustSetViewport(0, 0, 1, false)
	return page
}
