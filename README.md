# Todo CLI App (Go)

A simple command-line Todo application written in Go with a layered architecture:
- `main` for CLI handlers and menu flow
- `service` for business rules and validation
- `repository` for file-based persistence

## Features

- Create todos with title and description
- View:
  - all todos
  - pending todos
  - completed todos
- Update a todo
- Delete a todo (with confirmation)
- Toggle todo status (complete/incomplete)
- Persistent storage in `data/db.txt`

## Tech Stack

- Language: Go (`go 1.25.5` in `go.mod`)
- Storage: newline-delimited JSON records in a local text file

## Project Structure

```text
.
├── main.go              # CLI menu + handlers
├── models/
│   └── todo.go          # TODO model
├── service/
│   └── todo.go          # business logic + validation
├── repository/
│   └── todo.go          # data access layer
└── data/
    └── db.txt           # persisted todo data
```

## Getting Started

### Prerequisites

- Go installed (compatible with the version in `go.mod`)

### Run the app

```bash
go run main.go
```

### Build and run binary

```bash
go build -o todo-app .
./todo-app
```

## Menu Options

When started, the app shows:

1. Create New Todo
2. View All Todos
3. View Pending Todos
4. View Completed Todos
5. Update Todo
6. Delete Todo
7. Mark Todo Complete/Incomplete
8. Exit

## Validation Rules

- Title cannot be empty
- Title max length: 200 characters
- Description max length: 1000 characters
- Todo ID must be a positive integer for ID-based operations

## Data Format

Todos are stored in `data/db.txt` as one JSON object per line (newline-delimited JSON).  
Example record:

```json
{"ID":1,"Title":"Buy milk","Description":"2 liters","CreatedAt":"2026-02-21T10:00:00Z","Completed":false,"UpdatedAt":"2026-02-21T10:00:00Z"}
```

## Notes

- The app initializes the `data/` directory and `db.txt` automatically on startup.
- Repository operations are guarded with a mutex (`sync.RWMutex`) for safe concurrent access within the process.
