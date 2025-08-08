# Go Concurrency: Chat Application Project

## Project Overview

This project explores Go’s concurrency capabilities through a three-stage development process designed for progressively more complex applications.

The goal is to deliver a resilient, tested Go application that incorporates error handling, concurrency, client-server networking, and database interactions. Git version control is maintained throughout development.

---

## Project Stages

### Stage 1: Basic Concurrency

- Uses Go goroutines and channels to run multiple concurrent tasks.
- Terminal-based interface outputs concurrent task completions.
- Simple, foundational code illustrating Go’s concurrency model.

### Stage 2: Real-Time CLI Chat Application

- Implements a TCP chat server that handles multiple clients concurrently.
- Clients connect through a command-line interface with a raw terminal input mode.
- Supports user registration and authentication with password hashing (bcrypt).
- Stores user data securely in an SQLite database (Local to the machine running the server code).
- Maintains chat history that new clients receive upon connecting (history is deleted once the server disconnects).
- Supports user commands like `/help` and `/quit`.
- Protects against multiple logins with the same username.
- Provides basic cyber security measures for authentication.
  
### Stage 3: (Potential Future Work)

- Develop an application on top of the chat platform.
- Examples include a personal finance CLI tool or a multiplayer game.
- Enhance with file I/O, reporting, spreadsheet integration, data visualization, and a simple GUI.

---

## Project Structure

```plaintext
.
├── stage_one/
│   ├── go.mod
│   └── main.go          # Concurrency example with goroutines and channels
├── stage_two_chat_app/
│   ├── client/
│   │   ├── client.go    # CLI chat client with raw terminal input handling
│   │   ├── go.mod
│   │   └── go.sum
│   ├── server/
│   │   ├── auth.go      # Authentication logic using bcrypt
│   │   ├── db.go        # SQLite database setup and user management
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── main.go      # TCP chat server with concurrent client handling
```


# Getting Started

## Prerequisites 

### For Client Code

- Go 1.24 or newer

### For Server Code

- Go 1.24 or newer
- SQLite3 on your system (optional, the program creates a SQLite DB file if missing)
- MinGW-W64 x86_64 15.1.0 
    (gcc.exe (MinGW-W64 x86_64-msvcrt-posix-seh, built by Brecht Sanders, r4) 15.1.0 Copyright (C) 2025 Free Software Foundation, Inc. This is free software; see the source for copying conditions.  There is NO warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.)

## Stage 1: Running the Concurrency Example
1. Navigate to `stage_one` directory. 
2. Run the program:

```bash
#In the terminal
go run main.go
```

## Stage 2: Running the Chat Application

### Server

1. Navigate to `stage_two_chat_app/server`.

2. Download dependencies and build:

```bash
#In the terminal
go mod tidy
go build
```

## Run the Server

```bash
#In the terminal
./server
#or
go run .
```

The server will start listening on TCP port **6000**.

---

## Client

Open a new terminal and navigate to `stage_two_chat_app/client`.

### Download dependencies and build:

```bash
#In the terminal
go mod tidy
go build
```

### Run the Client

```bash
#In the terminal
./client
#or
go run .
```
## Usage

- When connecting, you will be prompted if you have an account.
- Register a new account or login with existing credentials.
- Use `/help` to see available commands.
- Use `/quit` to exit the chat.
- Chat messages are broadcast to all connected users.
- Chat history is loaded for the new user when they join.

---

## Testing & Error Handling

- Basic error handling throughout all stages (e.g., failed DB queries, invalid input).
- Authentication failure limits to **5 attempts** before disconnect.
- Username uniqueness enforced both in registration and login.
- Concurrent access to shared resources (`clients` map, chat history) guarded by **mutexes**.

---

## Limitations & Future Improvements

- Limited unit/integration testing due to time and skill constraints.
- Potential plans for a **GUI or web client interface** in Stage 3.
- Potential for more advanced features like **message encryption** 
---

## Authors & Version Control

- The project is maintained in this **GitHub repository** for version tracking. Only able to use a rudimentary amount of features with the current knowledge of ***Git*** and ***GitHub*** known by both authors
- The authors and ideas of this project are <b><i>Joseph Duff</i></b> and <b><i>Hector Marerro</i></b>


