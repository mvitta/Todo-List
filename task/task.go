package task

import (
	"encoding/json"
	"fmt"
)

type Task struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func NewTask(title string, id int) *Task {
	t := &Task{
		Id:     id,
		Title:  title,
		Status: false,
	}
	return t
}

func (t *Task) Json() []byte {

	theJson, errorJson := json.Marshal(t)
	if errorJson != nil {
		fmt.Println("Error: ", errorJson)
		return []byte{}
	}
	return theJson
}
