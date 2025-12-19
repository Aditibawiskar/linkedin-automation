package connect

import (
	"fmt"
	"linkedin-automation/internal/storage"
	"linkedin-automation/internal/stealth"
	"github.com/go-rod/rod"
)

func SendInvite(page *rod.Page) {
	// Find the Connect button (using generic text match for robustness)
	
	connectBtn, err := page.ElementR("button", "Connect")
	if err != nil {
		fmt.Println("No Connect button found on this user.")
		return
	}

	personName := "Unknown User" 
  if storage.IsInvited(personName) {
		fmt.Println("Skipping - already invited:", personName)
		return 
	}

	// After clicking success:
    storage.SaveInvite(personName)
	stealth.SleepRandom(1000, 3000) 
	connectBtn.MustClick()

	
	addNoteBtn, err := page.ElementR("button", "Add a note")
	if err == nil {
		stealth.SleepRandom(500, 1500)
		addNoteBtn.MustClick()
		
		// Type the message
		textArea := page.MustElement("textarea[name='message']")
		stealth.TypeLikeHuman(textArea, "Hi! I saw your profile and would love to connect.")
		
		// Click Send
		page.MustElementR("button", "Send").MustClick()
		fmt.Println("Invite sent with note!")
	}
}