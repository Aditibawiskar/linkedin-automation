package auth

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/go-rod/rod"
	"linkedin-automation/internal/stealth" // Ensure this path matches your module name
)

// AttemptLogin tries to log in using cookies first, then credentials
func AttemptLogin(browser *rod.Browser, page *rod.Page) {
	// --- STEP 1: TRY LOADING COOKIES ---
	fmt.Println("🤖 Bot: Checking for existing session cookies...")
	err := LoadCookies(browser, "cookies.json")
	if err == nil {
		fmt.Println("🤖 Bot: Cookies loaded. Verifying session...")
		page.MustNavigate("https://www.linkedin.com/feed/")
		page.MustWaitLoad()
		stealth.SleepRandom(2000, 3000)

		// Check if we are actually logged in
		if isLoggedIn(page) {
			fmt.Println("✅ Session Valid! Skipped Login.")
			return
		}
		fmt.Println("⚠️  Cookies expired or invalid. Proceeding to manual login...")
	}

	// --- STEP 2: MANUAL LOGIN (If cookies failed) ---
	email := os.Getenv("LINKEDIN_EMAIL")
	pass := os.Getenv("LINKEDIN_PASSWORD")

	if email == "" || pass == "" {
		log.Fatal("Error: LINKEDIN_EMAIL or LINKEDIN_PASSWORD not found in .env file")
	}

	fmt.Println("🤖 Bot: Navigating to Login Page...")
	page.MustNavigate("https://www.linkedin.com/login")
	page.MustWaitLoad()
	stealth.SleepRandom(2000, 4000)

	// Type Email
	fmt.Println("🤖 Bot: Typing Email...")
	emailField := page.MustElement("#username")
	stealth.TypeLikeHuman(emailField, email)
	stealth.SleepRandom(1000, 2000)

	// Type Password
	fmt.Println("🤖 Bot: Typing Password...")
	passField := page.MustElement("#password")
	stealth.TypeLikeHuman(passField, pass)
	stealth.SleepRandom(1000, 3000)

	// Click Sign In
	fmt.Println("🤖 Bot: Clicking Sign In...")
	page.MustElement("button[type='submit']").MustClick()

	// --- STEP 3: SECURITY CHECKPOINT & WAIT ---
	fmt.Println("🤖 Bot: Checking for security challenges...")

	// Wait logic: Check specifically for feed or security challenge
	waitForFeedOrChallenge(page)

	// --- STEP 4: SAVE NEW COOKIES ---
	fmt.Println("🤖 Bot: Login confirmed. Saving new cookies...")
	if err := SaveCookies(browser, "cookies.json"); err != nil {
		fmt.Printf("⚠️ Warning: Could not save cookies: %v\n", err)
	} else {
		fmt.Println("💾 Cookies saved to 'cookies.json' for next time.")
	}
}

func isLoggedIn(page *rod.Page) bool {

	found, _, _ := page.Timeout(5 * time.Second).Has(".global-nav__content")
	return found
}

// Helper to handle the "Wait for user to solve captcha" loop
func waitForFeedOrChallenge(page *rod.Page) {
	// Loop until we see the feed URL
	for {
		info, _ := page.Info()
		if info.URL == "https://www.linkedin.com/feed/" || info.URL == "https://www.linkedin.com/" {
			fmt.Println("✅ Feed detected! Resuming automation...")
			return
		}

		// If not on feed, warn user and wait
		fmt.Println("⏳ Waiting for Home Feed... (If Captcha appeared, please solve it manually)")
		time.Sleep(3 * time.Second)
	}
}
