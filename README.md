# Snippetbox

A web application for creating and viewing text snippets, built as an educational project following the "Let's Go!" book by Alex Edwards.

## 🎯 Project Overview

Snippetbox is a learning project that demonstrates web development concepts in Go. It's designed to teach fundamental web development principles including:

- HTTP routing and handlers
- Web server setup and configuration
- Request/response handling
- Basic web application architecture

## 🚀 Features

Currently implemented:
- **Home page** (`/`) - Displays a welcome message
- **Snippet view** (`/snippet/view`) - Placeholder for viewing specific snippets
- **Snippet creation** (`/snippet/create`) - Placeholder for creating new snippets

## 🛠️ Technology Stack

- **Language**: Go 1.24.5
- **Web Framework**: Standard library `net/http`
- **Routing**: `http.ServeMux`
- **Server**: Built-in HTTP server

## 📋 Prerequisites

- Go 1.24.5 or later
- Basic knowledge of Go syntax

## 🏃‍♂️ Getting Started

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd snippetbox
```

2. The project uses Go modules, so dependencies will be managed automatically.

### Running the Application

1. Start the web server:
```bash
go run ./cmd/web
```
or with command-line flags. For see flags use:
```bash
go run ./cmd/web -help
```

example:
```bash
go run ./cmd/web -addr=":4000"
```

Redirect the stdout and stderr streams on disk-files when starting application:
```bash
go run ./cmd/web >>./log/info.log 2>>./log/error.log
```

2. Open your browser and navigate to:
```
http://localhost:4000
```

You should see "Hello from Snippetbox" displayed.

### Available Routes

- `http://localhost:4000/` - Home page
- `http://localhost:4000/snippet/view` - Snippet view page
- `http://localhost:4000/snippet/create` - Snippet creation page

## 📁 Project Structure

```
snippetbox/
├── go.mod          # Go module definition
├── main.go         # Main application entry point
└── README.md       # This file
```

## 🔧 Development

This is an educational project following the "Let's Go!" book. The application is currently in its early stages and will be expanded with additional features as the learning progresses.

### Current Implementation

The application currently includes:
- Basic HTTP server setup on port 4000
- Three route handlers:
  - `home()` - Handles the root path
  - `snippetView()` - Handles snippet viewing
  - `snippetCreate()` - Handles snippet creation
- Simple request routing using `http.ServeMux`

## 📚 Learning Resources

This project is based on the book "Let's Go!" by Alex Edwards, which teaches web development with Go. The book covers:

- Building web applications with Go
- HTTP routing and middleware
- Database integration
- Security best practices
- Testing web applications

## 🤝 Contributing

This is an educational project, but suggestions and improvements are welcome! Feel free to:

- Report issues
- Suggest improvements
- Share your learning experience

## 📄 License

This project is created for educational purposes as part of learning Go web development.

## 🙏 Acknowledgments

- Alex Edwards for the excellent "Let's Go!" book
- The Go community for the robust standard library

---

**Note**: This is a work-in-progress educational project. Features and structure will evolve as the learning journey continues.
