package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"learn/models"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

var ErrNotFound = errors.New("todo not found")

type Repository struct {
	filePath string
	mu sync.RWMutex
}

type TodoRepository interface {
	Init() error
	SaveTodo(todo models.TODO) error
	GetAll() ([]models.TODO, error)
	GetOne(ID int) (models.TODO, error)
	UpdateOne(newTodo models.TODO, ID int) error
	DeleteOne(ID int) error
}

func New(filePath string) *Repository {
	return &Repository{
		filePath: filePath,
	}
}

func ( r *Repository) Init() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	if _, err := os.Stat(r.filePath);  os.IsNotExist(err) {
		file, err := os.Create(r.filePath)
		if err != nil {
			return fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()
	}

	return nil
}

func (r *Repository) SaveTodo(todo models.TODO) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todos, err := r.readAllTodos()
	if err != nil {
		return err
	}

	for _, existing := range todos {
		if existing.ID == todo.ID {
			return fmt.Errorf("todo with ID %d already exists", todo.ID)
		}
	}

	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(todo); err != nil {
		return fmt.Errorf("failed to encode todo: %w", err)
	}

	return nil
}


func (r *Repository) GetAll() ([]models.TODO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.readAllTodos()
}

func (r *Repository) GetOne(ID int) (models.TODO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos, err := r.readAllTodos()
	if err != nil {
		return models.TODO{}, err
	}

	for _, todo := range todos {
		if todo.ID == ID {
			return todo, nil
		}
	}

	return models.TODO{}, ErrNotFound
}

func (r *Repository) UpdateOne(newTodo models.TODO, ID int) (error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todos, err := r.readAllTodos()
	if err != nil {
		return fmt.Errorf("error reading the todos %w", err)
	}

	found := false
	for i, t := range todos {
		if t.ID == ID {
			newTodo.ID = ID
			todos[i] = newTodo
			found = true
			break;
		}
	}

		if !found {
		return ErrNotFound
	}

	return r.writeAllTodos(todos)
}

func (r *Repository) DeleteOne(ID int) (error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todos, err := r.readAllTodos()
	if err != nil {
		return err
	}

	originalLength := len(todos)

	for i, todo := range todos {
		if todo.ID == ID {
			todos = slices.Delete(todos, i, i+1)
			break
		}
	}

	if len(todos) == originalLength {
		return ErrNotFound
	}

	return r.writeAllTodos(todos)
}


func (r *Repository) readAllTodos() ([]models.TODO, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.TODO{}, nil
		}
	}

	defer file.Close()

	var todos []models.TODO

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
			var todo models.TODO
		if err := json.Unmarshal(scanner.Bytes(), &todo); err != nil {
			return nil, fmt.Errorf("failed to decode todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return todos, nil
}

func (r *Repository) writeAllTodos(todos []models.TODO) error {
	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, todo := range todos {
		if err := encoder.Encode(todo); err != nil {
			return fmt.Errorf("failed to encode todo: %w", err)
		}
	}

	return nil
}