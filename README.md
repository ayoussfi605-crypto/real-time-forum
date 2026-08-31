# Real-Time Forum

A real-time forum application featuring real-time private messaging, posts, comments, categories, and interactions. Built with a Go backend and a Vanilla JavaScript frontend.

## 🚀 Features

* **User Authentication**: Secure registration and login using bcrypt for password hashing and UUID for session management.
* **Forum Core**:
  * Create posts and categorize them.
  * Comment on posts.
  * Like and dislike posts and comments.
* **Real-Time Chat**:
  * Real-time private messaging powered by WebSockets.
  * Real-time online / offline user statuses.
  * Real-time typing indicators.
  * Unread message notifications and badges.
  * Chat history pagination and chronologically ordered messages.
* **Responsive UI**: Clean and intuitive user interface built with HTML, CSS, and Vanilla JavaScript.

## 🛠️ Tech Stack

* **Backend**: Go (Golang)
* **Database**: SQLite3 (`mattn/go-sqlite3`)
* **WebSockets**: Gorilla WebSocket (`gorilla/websocket`)
* **Frontend**: HTML5, CSS3, Vanilla JavaScript
* **Security**: `golang.org/x/crypto` for password hashing

## 📂 Project Structure

```
real-time-forum/
├── backend/                  # Go Backend Application
│   ├── database/             # SQLite database setup and queries
│   ├── handlers/             # HTTP endpoint handlers (Auth, Posts, Comments, etc.)
│   ├── helpers/              # Utility functions (JSON responses, Validation, Hashing)
│   ├── middlewares/          # HTTP middlewares (Authentication)
│   ├── routes/               # API route definitions
│   ├── websocket/            # WebSocket hub, clients, and real-time event handling
│   ├── main.go               # Backend entry point
│   ├── go.mod                # Go module dependencies
│   └── forum                 # Compiled binary (if built)
└── frontend/                 # Static Frontend Application
    ├── css/                  # Stylesheets
    ├── js/                   # Vanilla JavaScript logic (chat, state, auth, etc.)
    └── index.html            # Main SPA entry point
```

## ⚙️ How to Run

1. **Clone the repository** (if you haven't already).
2. **Navigate to the backend directory**:
   ```bash
   cd real-time-forum/backend
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Run the Go server**:
   ```bash
   go run main.go
   ```
5. **Access the application**:
   Open your browser and navigate to the local server address (usually `http://localhost:8080`).

## 👥 Authors

Created and maintained by:
* [@aelyoussef](https://github.com/ayoussfi605-crypto)
* [@mbenboua](https://github.com/mohamedben003)

