package auth

import (
	"encoding/json"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// LoadCookies loads the session from a JSON file to bypass login
func LoadCookies(page *rod.Page, filePath string) error {
	// 1. Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 2. Unmarshal into a temporary struct matching the JSON
	var cookies []*proto.NetworkCookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return err
	}

	// 3. Convert "NetworkCookie" to "NetworkCookieParam"
	// The library requires this specific conversion
	var params []*proto.NetworkCookieParam
	for _, c := range cookies {
		// We copy the fields we need
		p := &proto.NetworkCookieParam{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Secure:   c.Secure,
			HTTPOnly: c.HTTPOnly,
			SameSite: c.SameSite,
			Expires:  c.Expires,
		}
		params = append(params, p)
	}

	// 4. Set the cookies in the browser
	return page.SetCookies(params)
}