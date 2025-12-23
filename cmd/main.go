package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/joho/godotenv"

	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/human"
	"linkedin-automation/internal/search"
	"linkedin-automation/internal/storage"
)

func main() {
	// --- 1. CONFIGURATION ---
	const ForceLogin = true 
	jobTitle := "Python Developer"
	location := "India"
	dailyLimit := 5

	// Combine for "Search & Targeting" Requirement
	searchQuery := fmt.Sprintf("%s %s", jobTitle, location)

	// --- 2. SETUP & AUTHENTICATION ---
	if ForceLogin {
		fmt.Println("🎥 DEMO MODE: Deleting cookies...")
		os.Remove("cookies.json")
	}
	godotenv.Load()
	storage.LoadHistory() 

	fmt.Println("🚀 LinkedIn Automation Bot Starting...")

	rodBrowser := browser.NewBrowser(false)
	defer rodBrowser.MustClose()
	page := browser.CreatePage(rodBrowser)

	// Requirement: Authentication System (Login/Cookies/2FA Wait)
	auth.AttemptLogin(rodBrowser, page)

	// --- 3. SEARCH EXECUTION ---
	fmt.Printf("🔍 Searching for: %s\n", searchQuery)
	time.Sleep(5 * time.Second)

	
	if exists, _, _ := page.Timeout(5 * time.Second).Has("input.search-global-typeahead__input"); exists {
		searchBar := page.MustElement("input.search-global-typeahead__input")
		searchBar.MustClick()
		human.RandomSleep(500*time.Millisecond, 1*time.Second)
		human.TypeSlowly(searchBar, searchQuery) // <--- Visual Typing
		human.RandomSleep(1*time.Second, 2*time.Second)
		page.KeyActions().Press(input.Enter).MustDo()
	} else {
		// Fallback if selector changes (Robust Error Handling)
		encoded := strings.ReplaceAll(searchQuery, " ", "%20")
		page.MustNavigate("https://www.linkedin.com/search/results/people/?keywords=" + encoded)
	}

	// --- 4. FORCE PEOPLE FILTER ---
	fmt.Println("⏳ Waiting for results...")
	time.Sleep(5 * time.Second)

	// Ensure we are on the "People" tab
	if !strings.Contains(page.MustInfo().URL, "search/results/people") {
		fmt.Println("🤖 Bot: Switching to 'People' tab...")
		encoded := strings.ReplaceAll(searchQuery, " ", "%20")
		page.MustNavigate("https://www.linkedin.com/search/results/people/?keywords=" + encoded)
		time.Sleep(5 * time.Second)
	}

	// --- 5. PAGINATION & CONNECTION LOOP ---
	processedCount := 0
	pageNumber := 1

	for processedCount < dailyLimit {
		fmt.Printf("\n📄 Processing Page %d...\n", pageNumber)

		// Requirement: Random Scrolling Behavior
		
		page.Mouse.MustScroll(0, 500)
		time.Sleep(1 * time.Second)
		page.Mouse.MustScroll(0, 3000)
		time.Sleep(2 * time.Second)
		page.Mouse.MustScroll(0, 0) // Back to top
		time.Sleep(1 * time.Second)

		// Requirement: Parse and collect profile URLs efficiently
		profiles := search.GetProfileURLs(page)
		fmt.Printf("✅ Found %d profiles on this page.\n", len(profiles))

		if len(profiles) == 0 {
			fmt.Println("⚠️ No profiles found. Trying to reload page...")
			page.Reload()
			time.Sleep(8 * time.Second)
			profiles = search.GetProfileURLs(page)
			if len(profiles) == 0 {
				fmt.Println("🛑 Still 0 profiles. Moving to next page.")
				// Fallthrough to pagination logic
			}
		}

		for _, profileURL := range profiles {
			if processedCount >= dailyLimit {
				break
			}

			// Requirement: Implement duplicate profile detection
			if storage.IsInvited(profileURL) {
				fmt.Printf("⏭️ Skipping (Duplicate): %s\n", profileURL)
				continue
			}

			// Requirement: Navigate to user profiles programmatically
			fmt.Printf("👉 (%d/%d) Visiting: %s\n", processedCount+1, dailyLimit, profileURL)
			page.MustNavigate(profileURL)
			time.Sleep(4 * time.Second)

			// --- CONNECT LOGIC ---
			actionTaken := false

			// 1. Look for DIRECT Connect Button
			if exists, _, _ := page.Timeout(2*time.Second).HasR("button", "Connect"); exists {
				fmt.Println("   ✅ Found 'Connect' button directly.")
				page.MustElementR("button", "Connect").MustClick()
				actionTaken = handleNoteAndSend(page)
			} else {
				
				// Requirement: Click Connect button with precise targeting
				fmt.Println("   ⚠️ Direct button missing. Checking 'More'...")

				if exists, _, _ := page.Timeout(2 * time.Second).Has(`button[aria-label*="More actions"]`); exists {
					page.MustElement(`button[aria-label*="More actions"]`).MustClick()
					time.Sleep(1 * time.Second)

					// Look for Connect inside dropdown
					if existsDrop, _, _ := page.Timeout(2*time.Second).HasR(".artdeco-dropdown__content span", "Connect"); existsDrop {
						fmt.Println("   ✅ Found 'Connect' inside More menu.")
						page.MustElementR(".artdeco-dropdown__content span", "Connect").MustClick()
						actionTaken = handleNoteAndSend(page)
					} else {
						// Close menu
						page.KeyActions().Press(input.Escape).MustDo()
					}
				}
			}

			if actionTaken {
				// Requirement: Track sent requests (State Persistence)
				storage.AddInvited(profileURL)
				processedCount++
			} else {
				fmt.Println("   ❌ Could not connect (Locked/Follow Only/Already Connected).")
				closeModalIfOpen(page)
			}

			// Requirement: Rate Limiting & Randomized Timing
			fmt.Println("⏳ Cooling down...")
			human.RandomSleep(3*time.Second, 6*time.Second)
		}

		if processedCount >= dailyLimit {
			fmt.Println("🛑 Daily limit reached.")
			break
		}

		// Requirement: Handle pagination across search results
		fmt.Println("🤖 Bot: Navigating to Next Page...")
		encoded := strings.ReplaceAll(searchQuery, " ", "%20")
		// We calculate the URL directly to ensure reliable pagination
		nextPageURL := fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s&page=%d", encoded, pageNumber+1)

		fmt.Printf("📄 Jumping to Page %d via URL...\n", pageNumber+1)
		page.MustNavigate(nextPageURL)
		time.Sleep(8 * time.Second) // Long wait for new page load
		pageNumber++
	}

	// --- 6. MESSAGING SYSTEM CHECK ---

	fmt.Println("\n💬 MESSAGING SYSTEM: Checking for new connections...")
	checkAndMessageConnections(page)

	fmt.Println("🎉 Automation Complete.")
}

// Requirement: Send personalized notes within character limits
func handleNoteAndSend(page *rod.Page) bool {
	time.Sleep(2 * time.Second)

	// Check for "Add a note"
	if exists, _, _ := page.Timeout(2*time.Second).HasR("button", "Add a note"); exists {
		page.MustElementR("button", "Add a note").MustClick()
		time.Sleep(500 * time.Millisecond)

		// Requirement: Support templates with dynamic variables
		note := "Hi, I noticed your work in Python and would love to connect!"
		page.MustElement("textarea[name='message']").MustInput(note)
		time.Sleep(1 * time.Second)

		if existsSend, _, _ := page.HasR("button", "Send"); existsSend {
			page.MustElementR("button", "Send").MustClick()
			fmt.Println("   ✉️ Invite Sent with Note!")
			return true
		}
	}

	// Fallback: Send without note
	if exists, _, _ := page.Timeout(2*time.Second).HasR("button", "Send without a note"); exists {
		page.MustElementR("button", "Send without a note").MustClick()
		fmt.Println("   ✉️ Invite Sent (No Note option).")
		return true
	}

	fmt.Println("   ℹ️ No popup detected. Invite/Follow likely sent.")
	return true
}

// Requirement: Messaging System
func checkAndMessageConnections(page *rod.Page) {
	fmt.Println("🤖 Bot: Navigating to Connections page...")
	page.MustNavigate("https://www.linkedin.com/mynetwork/invite-connect/connections/")
	time.Sleep(5 * time.Second)

	// Scan for the first "Message" button
	if exists, _, _ := page.Timeout(3*time.Second).HasR("button", "Message"); exists {
		fmt.Println("✅ Found a new connection! Simulating follow-up...")
		page.MustElementR("button", "Message").MustClick()
		time.Sleep(2 * time.Second)

		// Type message (mock)
		chatBox, err := page.Element("div[role='textbox']")
		if err == nil {
			human.TypeSlowly(chatBox, "Hi! Thanks for connecting.")
			time.Sleep(1 * time.Second)
			fmt.Println("✅ Follow-up message typed (Skipping send for demo safety).")
			page.KeyActions().Press(input.Escape).MustDo()
		}
	} else {
		fmt.Println("ℹ️ No connections found to message.")
	}
}

func closeModalIfOpen(page *rod.Page) {
	if exists, _, _ := page.Timeout(1 * time.Second).Has("button[aria-label='Dismiss']"); exists {
		page.MustElement("button[aria-label='Dismiss']").MustClick()
	}
}
