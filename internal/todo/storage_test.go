package todo

import "testing"

type memStore struct {
    data []Task
}

func (m *memStore) Load() ([]Task, error) {
    return m.data, nil
}

func (m *memStore) Save(t []Task) error {
    m.data = t
    return nil
}

func TestAddAndComplete(t *testing.T) {
    svc, _ := NewService(&memStore{})
    task, _ := svc.Add("write tests")
    if task.ID != 0{
        t.Fatalf("expected id 0 but got %d", task.ID)
    }

    _ = svc.Complete(task.ID)
    list, _ := svc.List()
    if !list[0].Done {
        t.Fatal("task should be done")
    }
}
