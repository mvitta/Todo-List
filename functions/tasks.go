package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	t "main/task"
)

var Tasks = []t.Task{
	{Id: 1, Title: "Task 1", Status: false},
}

func SliceTasksToJson() []byte {
	sliceJson, errorJson := json.Marshal(Tasks)
	if errorJson != nil {
		fmt.Println("Error: ", errorJson)
	}
	return sliceJson
}

func FindTaskByID(id int) (*t.Task, error) {
	for i := 0; i < len(Tasks); i++ {
		task := Tasks[i]
		if task.Id == id {
			return &task, nil
		}
	}
	return nil, errors.New("the id does not exist")
}

func FindAndUpdateStatus(id int, status bool) error {
	for i := 0; i < len(Tasks); i++ {
		task := &Tasks[i]
		if task.Id == id {
			task.Status = status
			return nil
		}
	}
	return errors.New("the id does not exist")
}
