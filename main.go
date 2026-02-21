package main

import (
	"bufio"
	"fmt"
	"learn/models"
	"learn/repository"
	"learn/service"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	repo *repository.Repository
	svc  service.TodoServiceInterface
)

func init() {

	dataDir := "./data"

	dbFilePath := filepath.Join(dataDir, "db.txt")
	repo = repository.New(dbFilePath)

	if err := repo.Init(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	svc = service.New(repo)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		displayMainMenu()
		choice := getUserInput(reader, "Choose option: ")

		switch strings.TrimSpace(choice) {
		case "1":
			createTodoHandler(reader)
		case "2":
			viewAllTodosHandler()
		case "3":
			viewPendingTodosHandler()
		case "4":
			viewCompletedTodosHandler()
		case "5":
			updateTodoHandler(reader)
		case "6":
			deleteTodoHandler(reader)
		case "7":
			toggleTodoStatusHandler(reader)
		case "8":
			fmt.Println("Exiting Todo App. Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid option. Please try again.")
		}
		fmt.Println()
	}
}

func displayMainMenu() {
	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║        TODO APP - Main Menu         ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║ 1. Create New Todo                  ║")
	fmt.Println("║ 2. View All Todos                   ║")
	fmt.Println("║ 3. View Pending Todos               ║")
	fmt.Println("║ 4. View Completed Todos             ║")
	fmt.Println("║ 5. Update Todo                      ║")
	fmt.Println("║ 6. Delete Todo                      ║")
	fmt.Println("║ 7. Mark Todo Complete/Incomplete    ║")
	fmt.Println("║ 8. Exit                             ║")
	fmt.Println("╚════════════════════════════════════╝")
}

func createTodoHandler(reader *bufio.Reader) {
	fmt.Println("\n--- Create New Todo ---")

	title := getUserInput(reader, "Enter title (max 200 chars): ")
	description := getUserInput(reader, "Enter description (max 1000 chars): ")

	todo, err := svc.CreateTodo(strings.TrimSpace(title), strings.TrimSpace(description))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Todo created successfully!\n")
	fmt.Printf("   ID: %d | Title: %s\n", todo.ID, todo.Title)
}

func viewAllTodosHandler() {
	fmt.Println("\n--- All Todos ---")

	todos, err := svc.GetAllTodos()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No todos found. Create one to get started!")
		return
	}

	displayTodos(todos)
}

func viewPendingTodosHandler() {
	fmt.Println("\n--- Pending Todos ---")

	todos, err := svc.GetPendingTodos()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No pending todos. Great job!")
		return
	}

	displayTodos(todos)
}

func viewCompletedTodosHandler() {
	fmt.Println("\n--- Completed Todos ---")

	todos, err := svc.GetCompletedTodos()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No completed todos yet.")
		return
	}

	displayTodos(todos)
}

func updateTodoHandler(reader *bufio.Reader) {
	fmt.Println("\n--- Update Todo ---")

	todos, err := svc.GetAllTodos()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No todos to update.")
		return
	}

	displayTodos(todos)

	idStr := getUserInput(reader, "Enter todo ID to update: ")
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("❌ Invalid ID format.")
		return
	}

	// Verify todo exists
	_, err = svc.GetTodoByID(id)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	title := getUserInput(reader, "Enter new title: ")
	description := getUserInput(reader, "Enter new description: ")

	err = svc.UpdateTodo(id, strings.TrimSpace(title), strings.TrimSpace(description))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Todo updated successfully!\n")
}

func deleteTodoHandler(reader *bufio.Reader) {
	fmt.Println("\n--- Delete Todo ---")

	todos, err := svc.GetAllTodos()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No todos to delete.")
		return
	}

	displayTodos(todos)

	idStr := getUserInput(reader, "Enter todo ID to delete: ")
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("❌ Invalid ID format.")
		return
	}

	// Confirmation
	confirm := getUserInput(reader, "Are you sure? (yes/no): ")
	if strings.ToLower(strings.TrimSpace(confirm)) != "yes" {
		fmt.Println("❌ Deletion cancelled.")
		return
	}

	err = svc.DeleteTodo(id)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Todo deleted successfully!\n")
}

func toggleTodoStatusHandler(reader *bufio.Reader) {
	fmt.Println("\n--- Toggle Todo Status ---")

	todos, err := svc.GetAllTodos()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No todos to update.")
		return
	}

	displayTodos(todos)

	idStr := getUserInput(reader, "Enter todo ID to toggle: ")
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("❌ Invalid ID format.")
		return
	}

	todo, err := svc.GetTodoByID(id)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	var toggleErr error
	if todo.Completed {
		toggleErr = svc.MarkTodoAsIncomplete(id)
	} else {
		toggleErr = svc.MarkTodoAsCompleted(id)
	}

	if toggleErr != nil {
		fmt.Printf("❌ Error: %v\n", toggleErr)
		return
	}

	status := "incomplete"
	if !todo.Completed {
		status = "completed"
	}
	fmt.Printf("✅ Todo marked as %s!\n", status)
}

func displayTodos(todos []models.TODO) {
	fmt.Println("┌─────┬──────────────────────┬──────────────────────────┬─────────────┐")
	fmt.Println("│ ID  │ Title                │ Description              │ Status      │")
	fmt.Println("├─────┼──────────────────────┼──────────────────────────┼─────────────┤")

	for _, todo := range todos {
		title := todo.Title
		description := todo.Description

		// Truncate long strings
		if len(title) > 20 {
			title = title[:17] + "..."
		}
		if len(description) > 24 {
			description = description[:21] + "..."
		}

		status := "Pending"
		if todo.Completed {
			status = "✓ Completed"
		}

		fmt.Printf("│ %-3d │ %-20s │ %-24s │ %-11s │\n", todo.ID, title, description, status)
	}

	fmt.Println("└─────┴──────────────────────┴──────────────────────────┴─────────────┘")
}

func getUserInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
