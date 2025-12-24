package search

import (
	"fmt"
	"github.com/go-rod/rod"
	"strings"
)

func GetProfileURLs(page *rod.Page) []string {
	fmt.Println("   🕵️  Scanning for profile links (via JavaScript)...")

	results, err := page.Eval(`() => {
		return Array.from(document.querySelectorAll('a'))
			.map(a => a.href)
			.filter(href => href.includes('/in/'))
	}`)

	if err != nil {
		fmt.Println("   ❌ JS Error:", err)
		return []string{}
	}

	var urls []string
	seen := make(map[string]bool)

	// Convert the JS result to Go slice
	for _, val := range results.Value.Arr() {
		rawUrl := val.String()

		// Filter junk
		if strings.Contains(rawUrl, "/jobs/") ||
			strings.Contains(rawUrl, "linkedin.com/company/") ||
			strings.Contains(rawUrl, "linkedin.com/feed/") ||
			strings.Contains(rawUrl, "googletagmanager") {
			continue
		}

		// Clean URL
		if strings.Contains(rawUrl, "?") {
			parts := strings.Split(rawUrl, "?")
			rawUrl = parts[0]
		}

		if !seen[rawUrl] {
			urls = append(urls, rawUrl)
			seen[rawUrl] = true
		}
	}

	if len(urls) == 0 {
		fmt.Println("   ❌ Debug: Found 0 links.")
	} else {
		fmt.Printf("   ✅ Debug: Found %d valid links.\n", len(urls))
	}

	// Return up to 5
	if len(urls) > 5 {
		return urls[:5]
	}
	return urls
}
