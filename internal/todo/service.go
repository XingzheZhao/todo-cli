package todo

import (
	"errors"
	"sync"
)

// Service wraps Storage with domain logic. A tiny mutex protects
// the ID counter when the CLI is extended concurrent operations
type Service struct {
    s Storage
    mu sync.Mutex
    nextID int
}

// NewService loads tasks eagerly so nextID is correct.
func NewService(s Storage) (*Service, error) {
    tasks, err := s.Load()
    if err != nil {
        return nil, err
    }

    max := -1

    for _, t := range tasks {
        if t.ID > max{
            max = t.ID
        }
    }

    return &Service{s: s, nextID: max+1}, nil
}

// Add inserts a new task and persits the list
func (svc *Service) Add(text string) (Task, error) {
    svc.mu.Lock()
    defer svc.mu.Unlock()

    task :=  Task{ID: svc.nextID, Text: text, Done: false}
    svc.nextID++

    tasks, _ := svc.s.Load()
    tasks = append(tasks, task)
    if err := svc.s.Save(tasks); err != nil{
        return Task{}, err
    }
    return task, nil
}

// List returns all tasks
func (svc *Service) List() ([]Task, error) {return svc.s.Load()}

// Complete marks a task as done
func (svc *Service) Complete(id int) error{
    tasks, err := svc.s.Load()
    if err != nil {
        return err
    }

    found := false
    
    for i := range tasks{
        if tasks[i].ID == id{
            found = true
            tasks[i].Done = true
            break
        }
    }

    if !found{
        return errors.New("id not found")
    }
    return svc.s.Save(tasks)
}
