package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Title       string
	IsCompleted bool
}

var TodoList []Task

func main() {

	err := LoadTodoList()
	if err != nil {
		log.Fatal(err)
	}

	NewTaskptr := flag.Bool("n", false, "Enter a new Task")
	ListTaskptr := flag.Bool("l", false, "List All incomplete Tasks")
	CmpltTaskptr := flag.Bool("c", false, "Mark A Task As Complete")

	flag.Parse()

	if *NewTaskptr {
		newTask := Task{flag.Args()[0], false}
		err := SaveToList(newTask)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task Saved")
	}
	if *ListTaskptr {
		fmt.Println("Index\tTitle\t\tCompleted")
		for i := 0; i < len(TodoList); i++ {
			fmt.Println(i+1, "\t", TodoList[i].Title, "\t", TodoList[i].IsCompleted)
		}
	}

	if *CmpltTaskptr {
		// Todo later

	}

}

func LoadTodoList() error {
	if _, err := os.Stat("Todo.csv"); errors.Is(err, os.ErrNotExist) {
		f, _ := os.Create("Todo.csv")
		f.Close()
	}

	file, err := os.Open("Todo.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		Title := rec[0]
		Completed, _ := strconv.ParseBool(rec[1])
		NewTask := Task{Title, Completed}
		TodoList = append(TodoList, NewTask)
	}

	return nil
}

func SaveToList(newTask Task) error {
	file, err := os.OpenFile("Todo.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	defer func() {
		writer.Flush()
		err = writer.Error()
	}()

	NewTask := []string{newTask.Title, strconv.FormatBool(newTask.IsCompleted), time.Now().String()}
	err = writer.Write(NewTask)
	if err != nil {
		return err
	}

	TodoList = append(TodoList, newTask)

	return nil
}
