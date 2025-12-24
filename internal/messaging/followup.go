package messaging

import (
    "github.com/go-rod/rod"
)

func CheckAndMessage(page *rod.Page, profileURL string, followUpMsg string) {
    page.MustNavigate(profileURL)
    page.MustWaitStable()

    // If the button says "Message" instead of "Pending" or "Connect", they accepted.
    msgBtn, err := page.ElementR("button", "Message")
    if err == nil {
        msgBtn.MustClick()
        page.MustWaitStable()
        // Type and send
        editor := page.MustElement(".msg-form__contenteditable")
        editor.MustInput(followUpMsg)
        page.MustElement(".msg-form__send-button").MustClick()
    }
}