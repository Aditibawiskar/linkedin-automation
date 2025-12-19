package search

import (
	"fmt"
	"github.com/go-rod/rod"
)

func SearchPeople(page *rod.Page, keyword string) {
	// 1. Navigate to the Search Page
	url := fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s", keyword)
	page.MustNavigate(url)
	page.MustWaitLoad()
	page.MustElement(".reusable-search__result-container")
	buttons := page.MustElementsX("//button[contains(., 'Connect')]")
    
	fmt.Printf("Found %d people to connect with.\n", len(buttons))
}