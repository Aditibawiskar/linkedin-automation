package search

import (
	"fmt"
	"github.com/go-rod/rod"
	"strings"
	"time"
)

// GetProfileURLs scrapes ALL profile links from the page
func GetProfileURLs(page *rod.Page) []string {
	page.Mouse.MustScroll(0, 0) // Top
	time.Sleep(500 * time.Millisecond)
	page.Mouse.MustScroll(0, 2000) // Bottom
	time.Sleep(2 * time.Second)
	page.Mouse.MustScroll(0, 0) // Back up
	time.Sleep(1 * time.Second)

	var urls []string

	// 2. BROAD SELECTOR: Get every single <a> tag on the page
	elements, err := page.Elements("a")
	if err != nil {
		return urls
	}

	seen := make(map[string]bool)
	for _, el := range elements {
		link, err := el.Property("href")
		if err != nil {
			continue
		}
		urlStr := link.String()
		if strings.Contains(urlStr, "/in/") &&
			!strings.Contains(urlStr, "miniProfile") &&
			!strings.Contains(urlStr, "headless") &&
			len(urlStr) > 25 {
			cleanLink := strings.Split(urlStr, "?")[0]
			if !seen[cleanLink] {
				urls = append(urls, cleanLink)
				seen[cleanLink] = true
			}
		}
	}
	return urls
}

// NextPage attempts to find and click the "Next" button
func NextPage(page *rod.Page) bool {
	fmt.Println("   Checking for Next Page...")

	// Scroll to bottom where pagination is
	page.Mouse.MustScroll(0, 5000)
	time.Sleep(2 * time.Second)

	// Try multiple selectors for the "Next" button
	selectors := []string{
		"button[aria-label='Next']",
		".artdeco-pagination__button--next",
		"button span:text('Next')",
	}

	for _, sel := range selectors {
		btn, err := page.Element(sel)
		if err == nil {
			// Check if disabled
			if disabled, _ := btn.Attribute("disabled"); disabled != nil {
				return false
			}

			btn.MustClick()
			return true
		}
	}

	return false
}
