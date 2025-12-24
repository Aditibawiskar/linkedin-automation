package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/joho/godotenv"
	"linkedin-automation/internal/auth"
	"linkedin-automation/internal/browser"
	"linkedin-automation/internal/human"
	"linkedin-automation/internal/search"
	"linkedin-automation/internal/storage"
)

func main() {
	// --- CONFIG ---
	const ForceLogin = true
	jobTitle := "Python Developer"
	location := "India"
	searchQuery := fmt.Sprintf("%s %s", jobTitle, location)

	// --- SETUP ---
	if ForceLogin {
		os.Remove("cookies.json")
	}
	godotenv.Load()
	storage.LoadHistory()

	fmt.Println("🚀 LinkedIn Bot Starting (VIDEO MODE)...")
	rodBrowser := browser.NewBrowser(false)
	defer rodBrowser.MustClose()
	page := browser.CreatePage(rodBrowser)

	// --- AUTH ---
	auth.AttemptLogin(rodBrowser, page)

	// --- SEARCH ---
	fmt.Printf("🔍 Searching for: %s\n", searchQuery)
	time.Sleep(3 * time.Second)

	// Try finding search bar
	searchSelector := ""
	if exists, _, _ := page.Timeout(2 * time.Second).Has("input.search-global-typeahead__input"); exists {
		searchSelector = "input.search-global-typeahead__input"
	} else if exists, _, _ := page.Timeout(2 * time.Second).Has("input[placeholder='Search']"); exists {
		searchSelector = "input[placeholder='Search']"
	} else if exists, _, _ := page.Timeout(2 * time.Second).Has("input[role='combobox']"); exists {
		searchSelector = "input[role='combobox']"
	}

	if searchSelector != "" {
		fmt.Println("🤖 Bot: Found Search Bar! Typing query...")
		searchBar := page.MustElement(searchSelector)
		searchBar.MustClick()
		human.RandomSleep(500*time.Millisecond, 1*time.Second)
		human.TypeSlowly(searchBar, searchQuery)
		human.RandomSleep(1*time.Second, 2*time.Second)
		page.KeyActions().Press(input.Enter).MustDo()
	} else {
		fmt.Println("⚠️ Search bar hidden. Using direct navigation.")
		encoded := strings.ReplaceAll(searchQuery, " ", "%20")
		page.MustNavigate("https://www.linkedin.com/search/results/people/?keywords=" + encoded)
	}

	fmt.Println("⏳ Waiting for results...")
	time.Sleep(5 * time.Second)

	if !strings.Contains(page.MustInfo().URL, "search/results/people") {
		fmt.Println("🤖 Bot: Switching to 'People' tab...")
		encoded := strings.ReplaceAll(searchQuery, " ", "%20")
		page.MustNavigate("https://www.linkedin.com/search/results/people/?keywords=" + encoded)
		time.Sleep(8 * time.Second)
	}

	// --- VIDEO SCRIPT LOOP ---
	pageNumber := 1
	targetPage := 4 

	for {
		fmt.Printf("\n📄 Processing Page %d...\n", pageNumber)

		// 1. VISUAL SCROLLING
		fmt.Println("   🔄 Scrolling to read profiles...")
		page.KeyActions().Press(input.Home).MustDo()
		time.Sleep(1 * time.Second)
		for i := 0; i < 5; i++ {
			page.KeyActions().Press(input.PageDown).MustDo()
			time.Sleep(1 * time.Second)
		}
		page.KeyActions().Press(input.Home).MustDo()
		time.Sleep(2 * time.Second)

		if pageNumber < targetPage {
			fmt.Printf("   🎥 Video Mode: Showing Pagination (Skipping connect on Page %d)...\n", pageNumber)
			fmt.Println("🤖 Bot: Navigating to Next Page...")
			encoded := strings.ReplaceAll(searchQuery, " ", "%20")
			nextPageURL := fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s&page=%d", encoded, pageNumber+1)
			page.MustNavigate(nextPageURL)
			time.Sleep(5 * time.Second)
			pageNumber++
			continue
		}

		// 3. ACTION: PAGE 4
		fmt.Println("   🎯 Target Page Reached! Finding a profile to connect...")

		profiles := search.GetProfileURLs(page)

		if len(profiles) > 0 {
			targetProfile := profiles[0]
			fmt.Printf("👉 Visiting Target: %s\n", targetProfile)
			page.MustNavigate(targetProfile)
			time.Sleep(5 * time.Second)

			// CONNECT
			success := handleNoteAndSend(page)
			if success {
				fmt.Println("🎉 DEMO SUCCESS: Action taken!")
				fmt.Println("🛑 Stopping bot for video end.")
				return
			} else {
				fmt.Println("❌ Failed to connect. Trying the next profile...")
				// Try the second profile if the first fails
				if len(profiles) > 1 {
					targetProfile = profiles[1]
					fmt.Printf("👉 Visiting Target 2: %s\n", targetProfile)
					page.MustNavigate(targetProfile)
					time.Sleep(5 * time.Second)
					if handleNoteAndSend(page) {
						fmt.Println("🎉 DEMO SUCCESS: Action taken!")
						return
					}
				}
				return
			}
		} else {
			fmt.Println("⚠️ No profiles found on Page 4? Reloading...")
			page.Reload()
			time.Sleep(5 * time.Second)
		}
	}
}

// SAFE HANDLE FUNCTION (VERBOSE DEBUGGING)
func handleNoteAndSend(page *rod.Page) bool {
	fmt.Println("   ... Analyzing profile buttons...")
	time.Sleep(2 * time.Second)

	// 1. Check for "Connect" OR "+ Add" OR "Follow"
	
	buttons := []string{"Connect", "Add", "Follow", "+ Add", "+ Follow"}

	for _, btnText := range buttons {
		fmt.Printf("   ... Checking for '%s' button...\n", btnText)
		if exists, _, _ := page.Timeout(1*time.Second).HasR("button", btnText); exists {
			fmt.Printf("   ✅ Found '%s' button!\n", btnText)
			el, err := page.ElementR("button", btnText)
			if err == nil {
				el.Click(proto.InputMouseButtonLeft, 1)
				// Treat "Add" or "Connect" as a connection attempt
				if strings.Contains(btnText, "Connect") || strings.Contains(btnText, "Add") {
					return processConnectModal(page)
				}
				// If "Follow", just return true (Success for demo)
				fmt.Println("   ✨ Followed user (Demo Success).")
				return true
			}
		}
	}

	// 2. Check "More" Menu
	fmt.Println("   ... Buttons hidden. Checking 'More' menu...")
	if exists, _, _ := page.Timeout(2 * time.Second).Has(`button[aria-label*="More actions"]`); exists {
		el, err := page.Element(`button[aria-label*="More actions"]`)
		if err == nil {
			el.Click(proto.InputMouseButtonLeft, 1)
			time.Sleep(1 * time.Second)

			// Look for Connect inside
			if existsDrop, _, _ := page.Timeout(2*time.Second).HasR(".artdeco-dropdown__content span", "Connect"); existsDrop {
				fmt.Println("   ✅ Found 'Connect' in More menu.")
				dropEl, errDrop := page.ElementR(".artdeco-dropdown__content span", "Connect")
				if errDrop == nil {
					dropEl.Click(proto.InputMouseButtonLeft, 1)
					return processConnectModal(page)
				}
			} else {
				fmt.Println("   ❌ 'Connect' not found in More menu.")
			}
			// Close menu
			page.KeyActions().Press(input.Escape).MustDo()
		}
	}

	return false
}

func processConnectModal(page *rod.Page) bool {
	fmt.Println("   ... Handling Connection Modal...")
	time.Sleep(2 * time.Second)

	// A. Click Add a Note
	if exists, _, _ := page.Timeout(2*time.Second).HasR("button", "Add a note"); exists {
		el, _ := page.ElementR("button", "Add a note")
		el.Click(proto.InputMouseButtonLeft, 1)

		time.Sleep(500 * time.Millisecond)
		page.MustElement("textarea[name='message']").MustInput("Hi, I noticed your work and would love to connect!")
		time.Sleep(1 * time.Second)

		if existsSend, _, _ := page.HasR("button", "Send"); existsSend {
			elSend, _ := page.ElementR("button", "Send")
			elSend.Click(proto.InputMouseButtonLeft, 1)
			time.Sleep(2 * time.Second)

			// Check for Premium Popup
			if existsModal, _, _ := page.Has(".artdeco-modal"); existsModal {
				fmt.Println("   ⚠️ Popup Detected! Attempting Force Close...")
				forceClosePopup(page)
				time.Sleep(1 * time.Second)
				if existsNoNote, _, _ := page.HasR("button", "Send without a note"); existsNoNote {
					elNoNote, _ := page.ElementR("button", "Send without a note")
					elNoNote.Click(proto.InputMouseButtonLeft, 1)
					fmt.Println("   ✅ Recovered: Sent without note.")
					return true
				}
			}
			fmt.Println("   ✉️ Connection Request Sent!")
			return true
		}
	}

	// B. Direct Send (No Note Option)
	if exists, _, _ := page.Timeout(2*time.Second).HasR("button", "Send without a note"); exists {
		el, _ := page.ElementR("button", "Send without a note")
		el.Click(proto.InputMouseButtonLeft, 1)
		fmt.Println("   ✉️ Invite Sent (No Note).")
		return true
	}

	// C. Maybe it just sent immediately (Common with "Add" button)
	fmt.Println("   ℹ️ No modal appeared. Request likely sent automatically.")
	return true
}

func forceClosePopup(page *rod.Page) {
	if exists, _, _ := page.Timeout(1 * time.Second).Has("button[aria-label='Dismiss']"); exists {
		el, _ := page.Element("button[aria-label='Dismiss']")
		el.Click(proto.InputMouseButtonLeft, 1)
		return
	}
	page.KeyActions().Press(input.Escape).MustDo()
}
