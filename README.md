# LinkedIn Automation – Internship Assignment (SubSpace)

## Overview
This project is a technical proof-of-concept built for the SubSpace Software Development Internship (Sep-Dec Batch). 

It implements a **Golang-based browser automation tool** that navigates LinkedIn, searches for specific job titles, and automates connection requests with personalized notes. It utilizes the **Rod** library with advanced stealth techniques to simulate human behavior.

⚠️ **Educational Purpose Only:** This tool is a proof-of-concept and strictly adheres to ethical automation standards.

---

## 🚀 Features Implemented
* **Stealth Browser Automation:** Uses `go-rod/stealth` to mask automation signals (webdriver flags, user-agent).
* **Human Simulation:** Implements Bézier curve mouse movements, randomized typing speeds, and variable delays.
* **Smart Login Detection:** Supports a "Hybrid Login" workflow where the bot pauses for manual 2FA/CAPTCHA entry, then automatically resumes automation.
* **Search & Targeting:** Automates searching for keywords (e.g., "Software Engineer") and parsing results.
* **Connection & Messaging:** Automatically clicks "Connect," adds a personalized note, and sends the request.
* **State Persistence:** Uses a local JSON database (`history.json`) to track invited users and prevent duplicate requests.

---

## 🛠️ Tech Stack
* **Language:** Go (Golang)
* **Library:** Rod (DevTools Protocol)
* **Stealth:** go-rod/stealth
* **Storage:** JSON (local persistence)

---

## 📂 Project Structure
```text
linkedin-automation/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── connect/             # Logic for sending invites & notes
│   ├── search/              # Logic for searching & filtering
│   ├── stealth/             # Human-like mouse & behavior utilities
│   └── storage/             # JSON-based state persistence
├── .env.example             # Config template
├── .gitignore               # Ignored files
├── go.mod                   # Go module definition
└── README.md                # Documentation

---

⚙️ Setup & Usage
Prerequisites
• Go 1.20+ installed
• Google Chrome installed

Installation
1. Clone the repository:

• git clone [https://github.com/Aditibawiskar/linkedin-automation.git](https://github.com/Aditibawiskar/  linkedin-automation.git)
cd linkedin-automation


2. Install dependencies:
```Bash
• go mod tidy

Running the Bot
1. Run the application:
```Bash
• go run cmd/main.go


2. Manual Login Step:

• The browser will launch in full screen.

• Action Required: Manually enter your email/password and solve any CAPTCHAs.

• Once you reach the LinkedIn "Home Feed," the bot will detect the login success and automatically take over.


3. Watch: The bot will search, scroll, and send invites automatically.

---

🛡️ Anti-Detection Strategy
To meet the assignment's stealth requirements, this tool implements:

• Randomized Viewport: Mimics standard laptop screen resolutions.

• Mouse Pathing: No straight-line movements; uses randomized curvature.

• Variable Timing: Actions are spaced by random intervals (e.g., 2s–5s) to mimic human "think time."

---

Author
Aditi Bawiskar Software Development Internship Applicant