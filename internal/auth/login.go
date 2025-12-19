package auth

import (
	"log"

	"github.com/go-rod/rod"

	"linkedin-automation/internal/stealth"
)

// OpenLoginPage navigates to LinkedIn login and demonstrates human interaction
func OpenLoginPage(page *rod.Page) error {
	log.Println("Navigating to LinkedIn login page")

	page.MustNavigate("https://www.linkedin.com/login")
	page.MustWaitLoad()

	stealth.SleepRandom(2000, 4000)
	stealth.RandomScroll(page)

	// Locate email field
	emailInput, err := page.Element(`input[name="session_key"]`)
	if err != nil {
		return err
	}

	stealth.HoverElement(emailInput)
	stealth.TypeLikeHuman(emailInput, "demo@example.com")

	stealth.SleepRandom(1500, 2500)

	log.Println("Login page interaction completed (demo only)")
	return nil
}
