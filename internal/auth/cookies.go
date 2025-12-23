package auth

import (
	"encoding/json"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// SaveCookies exports current cookies to a file
func SaveCookies(browser *rod.Browser, filePath string) error {
	cookies, err := browser.GetCookies()
	if err != nil {
		return err
	}
	data, err := json.Marshal(cookies)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// LoadCookies imports cookies from a file if it exists
func LoadCookies(browser *rod.Browser, filePath string) error {
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil // No cookies yet, proceed to normal login
	}
	if err != nil {
		return err
	}

	var cookies []*proto.NetworkCookieParam
	if err := json.Unmarshal(data, &cookies); err != nil {
		return err
	}
	return browser.SetCookies(cookies)
}