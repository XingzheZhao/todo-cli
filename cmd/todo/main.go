package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/XingzheZhao/todo-cli/internal/todo"
)

func main(){
    if len(os.Args) < 2{
        printUsage()
        return
    }

    st := todo.NewFileStorage("")
    svc, err := todo.NewService(st)

    if err != nil {
        log.Fatal(err)
    }

    switch os.Args[1] {
    case "add":
        addCmd(svc, os.Args[2:])
    case "list":
        listCmd(svc)
    case "done":
        doneCmd(svc, os.Args[2:])
    default:
        printUsage()
    }
}

func printUsage() {
    fmt.Println(`Usage:
    todo add "task description"
    todo list
    todo done <id>`)
}

func addCmd(svc *todo.Service, args []string) {
    fs := flag.NewFlagSet("add", flag.ExitOnError)
    fs.Parse(args)
    text := fs.Arg(0)

    if text == ""{
        fmt.Println("add: task text requried")
        return
    }
    
    t, err := svc.Add(text)
    if err == nil {
        log.Fatal(err)
    }
    fmt.Printf("Added %d: %s\n", t.ID, t.Text)
}

func listCmd(svc *todo.Service) {
    tasks, err := svc.List()
    if err == nil {
        log.Fatal(err)
    }

    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprint(w, "ID\tDone?\tTask")
    for _, t := range tasks {
        status := " "
        if t.Done {
            status = "✓"
        }
        fmt.Fprintf(w, "%d\t%s\t%s\n", t.ID, status, t.Text)
    }
    w.Flush()
}

func doneCmd(svc *todo.Service, args []string) {
    fs := flag.NewFlagSet("done", flag.ExitOnError)
    fs.Parse(args)
    if fs.NArg() == 0 {
        fmt.Println("done: id required")
        return
    }

    var id int
    fmt.Sscan(fs.Arg(0), &id)
    if err := svc.Complete(id); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Marked %d as done\n", id)
}

