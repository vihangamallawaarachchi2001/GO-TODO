package service

import (
	"errors"
	"fmt"
	"learn/models"
	"learn/repository"
	"time"
)

type TodoService struct {
	repo repository.TodoRepository
}

type TodoServiceInterface interface {
	CreateTodo(title, description string) (models.TODO, error)
	GetAllTodos() ([]models.TODO, error)
	GetTodoByID(ID int) (models.TODO, error)
	UpdateTodo(ID int, title, description string) error
	DeleteTodo(ID int) error
	MarkTodoAsCompleted(ID int) error
	MarkTodoAsIncomplete(ID int) error
	GetCompletedTodos() ([]models.TODO, error)
	GetPendingTodos() ([]models.TODO, error)
	GetNextTodoID() (int, error)
}

func New(repo repository.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) CreateTodo(title string, description string) (models.TODO , error) {
	if title == "" {
		return models.TODO{}, errors.New("title cannot be empty")
	}

	if len(title) > 200 {
		return models.TODO{}, errors.New("title cannot exceed 200 characters")
	}

	if len(description) > 1000 {
		return models.TODO{}, errors.New("description cannot exceed 1000 characters")
	}

	nextID, err := s.GetNextTodoID()

		if err != nil {
		return models.TODO{}, err
	}

	todo := models.TODO{
		ID:          nextID,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Completed:   false,
	}

	if err := s.repo.SaveTodo(todo); err != nil {
		return models.TODO{}, fmt.Errorf("Failed to save todo: %w", err)
	}

	return todo, nil
}

func (s *TodoService) GetAllTodos() ([]models.TODO, error) {
	todos, err := s.repo.GetAll()
	if err !=nil {
		return nil, fmt.Errorf("failed to retrieve todos: %w", err)
	}

	if todos == nil {
		return []models.TODO{}, nil
	}

	return todos, nil
}

func (s *TodoService) GetTodoByID(ID int) (models.TODO, error) {
	if ID <= 0 {
		return models.TODO{}, errors.New("invalid todo ID")
	}

	todo, err := s.repo.GetOne(ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return models.TODO{}, fmt.Errorf("todo with ID %d not found", ID)
		}
		return models.TODO{}, fmt.Errorf("failed to retrieve todo: %w", err)
	}

	return todo, nil
}

func (s *TodoService) UpdateTodo(ID int, title, description string) error {
	if ID <= 0 {
		return errors.New("invalid todo ID")
	}

	if title == "" {
		return errors.New("title cannot be empty")
	}

	if len(title) > 200 {
		return errors.New("title cannot exceed 200 characters")
	}

	if len(description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	todo, err := s.repo.GetOne(ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("todo with ID %d not found", ID)
		}
		return fmt.Errorf("failed to retrieve todo: %w", err)
	}

	todo.Title = title
	todo.Description = description
	todo.UpdatedAt = time.Now()

	if err := s.repo.UpdateOne(todo, ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("todo with ID %d not found", ID)
		}
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

func (s *TodoService) DeleteTodo(ID int) error {
	if ID <= 0 {
		return errors.New("invalid todo ID")
	}

	err := s.repo.DeleteOne(ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("todo with ID %d not found", ID)
		}
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}

func (s *TodoService) MarkTodoAsCompleted(ID int) error {
	if ID <= 0 {
		return errors.New("invalid todo ID")
	}

	todo, err := s.repo.GetOne(ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("todo with ID %d not found", ID)
		}
		return fmt.Errorf("failed to retrieve todo: %w", err)
	}

	if todo.Completed {
		return errors.New("todo is already completed")
	}

	todo.Completed = true
	todo.UpdatedAt = time.Now()

	if err := s.repo.UpdateOne(todo, ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("todo with ID %d not found", ID)
		}
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

func (s *TodoService) MarkTodoAsIncomplete(ID int) error {
	if ID <= 0 {
		return errors.New("invalid todo ID")
	}

	todo, err := s.repo.GetOne(ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("todo with ID %d not found", ID)
		}
		return fmt.Errorf("failed to retrieve todo: %w", err)
	}

	if !todo.Completed {
		return errors.New("todo is already incomplete")
	}

	todo.Completed = false
	todo.UpdatedAt = time.Now()

	if err := s.repo.UpdateOne(todo, ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("todo with ID %d not found", ID)
		}
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

func (s *TodoService) GetCompletedTodos() ([]models.TODO, error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve todos: %w", err)
	}

	var completed []models.TODO
	for _, todo := range todos {
		if todo.Completed {
			completed = append(completed, todo)
		}
	}

	return completed, nil
}

func (s *TodoService) GetPendingTodos() ([]models.TODO, error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve todos: %w", err)
	}

	var pending []models.TODO
	for _, todo := range todos {
		if !todo.Completed {
			pending = append(pending, todo)
		}
	}

	return pending, nil
}

func (s *TodoService) GetNextTodoID() (int, error) {
	todos, err := s.repo.GetAll()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve todos: %w", err)
	}

	maxID := 0
	for _, todo := range todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}

	return maxID + 1, nil
}
