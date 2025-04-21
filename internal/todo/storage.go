package todo

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
)


type Storage interface {
    Load() ([]Task, error)
    Save([]Task) error
}

type fileStorage struct{
    path string
}

func NewFileStorage(file string) Storage{
    if file == ""{
        home, _ := os.UserHomeDir()
        file = filepath.Join(home, ".todo.json")
    }
    return &fileStorage{path: file}
}

func (fs *fileStorage) Load() ([]Task, error) {
    if _, err := os.Stat(fs.path); errors.Is(err, os.ErrNotExist) {
        return []Task{}, nil
    }
    raw, err := os.ReadFile(fs.path)
    if err != nil {
        return nil, err
    }

    var tasks []Task
    if err := json.Unmarshal(raw, &tasks); err != nil{
        return nil, err
    } 
    return tasks, nil
}

func (fs *fileStorage) Save(tasks []Task) error {
    raw, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return err
    }
    // 0600 user-only permissions
    return os.WriteFile(fs.path, raw, 0o600)
}
