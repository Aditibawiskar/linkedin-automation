package auth

import (
	"fmt"
	"os"
	"time"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"linkedin-automation/internal/human"
)

func AttemptLogin(browser *rod.Browser, page *rod.Page) {

	if err := LoadCookies(browser, "cookies.json"); err == nil {
		fmt.Println("🤖 Bot: Cookies loaded. Verifying session...")
		page.MustNavigate("https://www.linkedin.com/feed/")
		time.Sleep(5 * time.Second)

		if page.MustInfo().URL == "https://www.linkedin.com/feed/" {
			fmt.Println("✅ Session valid. Skipped Login.")
			return
		}
	}

	// 2. Start Login Process
	fmt.Println("🤖 Bot: Navigating to Login Page...")
	page.MustNavigate("https://www.linkedin.com/login")
	time.Sleep(3 * time.Second)

	email := os.Getenv("LINKEDIN_EMAIL")
	pass := os.Getenv("LINKEDIN_PASSWORD")

	// 3. Human-like Typing
	fmt.Println("🤖 Bot: Typing Email...")
	if exists, _, _ := page.Has("#username"); exists {
		emailInput := page.MustElement("#username")
		human.TypeSlowly(emailInput, email)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("🤖 Bot: Typing Password...")
	if exists, _, _ := page.Has("#password"); exists {
		passInput := page.MustElement("#password")
		human.TypeSlowly(passInput, pass)
		time.Sleep(1 * time.Second)
	}

	// 4. Click Sign In
	fmt.Println("🤖 Bot: Clicking Sign In...")
	page.KeyActions().Press(input.Enter).MustDo()
	
	// 🛑 THE SECURITY CHECK PAUSE 🛑
	checkForSecurityCheck(page)

	// 5. Save Cookies
	fmt.Println("🤖 Bot: Login confirmed. Saving new cookies...")
	SaveCookies(browser, "cookies.json")
}

// This function freezes the bot until it sees the "Feed" (Home page)
func checkForSecurityCheck(page *rod.Page) {
	fmt.Println("⏳ CHECKING FOR CAPTCHA/SECURITY CHECK...")
	fmt.Println("👉 If you see a puzzle, please solve it manually now!")
	
	// Wait up to 120 seconds (2 minutes)
	for i := 0; i < 24; i++ { 
		url := page.MustInfo().URL
		if url == "https://www.linkedin.com/feed/" || url == "https://www.linkedin.com/" {
			fmt.Println("✅ Feed detected! Resuming automation...")
			return
		}
		
		// Still waiting...
		fmt.Printf("   ... Waiting for Feed (%d/120s)\n", (i+1)*5)
		time.Sleep(5 * time.Second)
	}
	
	fmt.Println("⚠️ Wait time over. Assuming login was successful or failed.")
}