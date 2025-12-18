# LinkedIn Automation – Proof of Concept (SubSpace Assignment)

## Overview

This project is a technical proof-of-concept built as part of the SubSpace Software Development Internship assignment.

It demonstrates browser automation fundamentals, human-like interaction patterns, stealth techniques, and clean Go architecture using the Rod library.

The goal of this assignment is NOT to automate LinkedIn for real-world usage, but to showcase:

- Browser automation fundamentals
- Anti-detection and stealth techniques
- Modular and maintainable Go code
- Understanding of real-world automation limitations

⚠️ Educational purpose only. This project must not be used in production.

---

## Tech Stack

- Language: Go
- Browser Automation: Rod (Chrome DevTools Protocol)
- Stealth Utilities: go-rod/stealth
- Browser: Google Chrome (headful mode)
- Operating System: Windows

---

## Features Implemented

- Automated Chrome launch using Rod
- Human-like mouse movement with randomized paths
- Randomized delays to simulate real user behavior
- Natural scrolling behavior
- Stealth browser configuration
- Graceful error handling and logging
- Modular folder structure following Go best practices

---

## Project Structure

linkedin-automation/
├── cmd/
│   └── main.go
├── internal/
│   ├── auth/
│   │   └── login.go
│   ├── stealth/
│   │   ├── mouse.go
│   │   ├── scroll.go
│   │   └── timing.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md

---

## Authentication & Login Limitation

OAuth-based authentication flows such as “Sign in with Google” intentionally block browser automation for security reasons.

This project demonstrates automation up to the authentication boundary, focusing on:

- Browser control
- Stealth techniques
- Human-like interaction patterns

Authenticated flows can be supported via session cookie reuse, which is documented conceptually but intentionally not enabled in this demo to respect platform security constraints.

---

## Environment Configuration

This project uses environment variables for configuration.

- .env – Local configuration file (not committed to version control)
- .env.example – Template describing required environment variables

---

## How to Run

Prerequisites:

- Go 1.20 or later
- Google Chrome installed
- Windows operating system

Steps:

1. Clone the repository
2. Copy .env.example to .env
3. Update values in .env if required
4. Run the application:

   go run cmd/main.go

The browser will launch and navigate to the LinkedIn login page, demonstrating human-like interaction behavior.

---

## Disclaimer

- This project does not attempt to bypass LinkedIn or Google security mechanisms.
- It must not be used for real account automation.
- Built strictly for educational and technical evaluation purposes.

---

## Author

Aditi Bawiskar  
Software Development Internship Applicant – SubSpace
