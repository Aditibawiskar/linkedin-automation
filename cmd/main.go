package main

import (
	"fmt"
	"time"
	
	// Import your internal packages
	"linkedin-automation/internal/search"
	"linkedin-automation/internal/connect"
	"linkedin-automation/internal/stealth"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	rodstealth "github.com/go-rod/stealth"
)

func main() {
	// 1. Configure the Browser
	u := launcher.New().
		Leakless(false).            // Fixes the "virus" error
		Headless(false).            // We want to see the browser
		// Set("window-size") is REMOVED so we can be full screen
		MustLaunch()

	// 2. Connect to the Browser
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// 3. Create the Page and MAXIMIZE it immediately
	// 3. Create the Page
	page := rodstealth.MustPage(browser)
    
    // 👇 ADD THIS LINE. It forces the internal view to be 1920x1080
    page.MustSetViewport(1540, 750, 1.0, false)
    
	page.MustWindowMaximize()

	// 4. Navigate to Login
	fmt.Println("Opening LinkedIn... Please log in manually!")
	page.MustNavigate("https://www.linkedin.com/login")

	// 5. WAIT FOR MANUAL LOGIN
	// The bot will pause here until it detects you have reached the Home Feed
	// It checks for the search bar or nav bar to know you are inside.
	fmt.Println("Waiting for you to sign in...")
	page.Race().Element(".global-nav__content").MustHandle(func(e *rod.Element) {
		fmt.Println("Login detected! Starting automation...")
	}).MustDo()

	// --- AUTOMATION STARTS HERE (After you log in) ---

	// 6. Search for People
	// Make sure you created internal/search/engine.go from my previous message!
	search.SearchPeople(page, "Software Engineer")
	
	// Human pause
	stealth.SleepRandom(2000, 4000)

	// 7. Send Connect Invites
	// Make sure you created internal/connect/invite.go from my previous message!
	connect.SendInvite(page)

	// 8. Keep window open for video recording
	fmt.Println("Done! Press Ctrl+C in the terminal to close.")
	time.Sleep(10 * time.Minute)
}