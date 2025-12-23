# LinkedIn Automation Bot (PoC)

A sophisticated Go-based automation tool utilizing `go-rod` to simulate human interactions on LinkedIn. This project demonstrates advanced browser automation, stealth techniques, and modular architecture.

## ⚠️ Disclaimer
**Educational Purpose Only.** This tool is a Proof-of-Concept designed for an internship assignment. Automating LinkedIn violates their Terms of Service.

## 🚀 Features
* **Stealth Automation:** Uses `go-rod/stealth` and randomized delays to mask bot behavior.
* **Smart Authentication:** Supports cookie persistence and pauses for manual 2FA/CAPTCHA handling.
* **Targeted Outreach:** Searches by Job Title + Location and filters results.
* **Safe Connection Logic:** Detects "Connect" vs "Message" buttons and handles "More" menus.
* **Duplicate Detection:** Maintains a `history.json` database to prevent spamming users twice.
* **Pagination:** Automatically traverses search result pages.

## 🛠️ Tech Stack
* Go (Golang)
* go-rod (Browser Automation)
* godotenv (Config management)

## 📦 Setup & Usage

1.  **Clone the repository**
    ```bash
    git clone <your-repo-url>
    cd linkedin-automation
    ```

2.  **Install Dependencies**
    ```bash
    go mod tidy
    ```

3.  **Configure Environment**
    Rename `.env.example` to `.env` and add your credentials:
    ```ini
    LINKEDIN_EMAIL=myuser@gmail.com
    LINKEDIN_PASSWORD=mypassword
    ```

4.  **Run the Bot**
    ```bash
    go run cmd/main.go
    ```

## 📂 Project Structure
* `cmd/main.go`: Entry point and main orchestration logic.
* `internal/auth`: Login and Cookie management.
* `internal/browser`: Rod browser setup and Stealth injection.
* `internal/search`: Profile parsing and pagination logic.
* `internal/storage`: JSON-based state persistence (History).
* `internal/human`: Human-like typing and sleep simulation.