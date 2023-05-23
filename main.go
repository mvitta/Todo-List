package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	t "main/task"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type RequestBody struct {
	Title string `json:"title"`
}

var tasks = []t.Task{
	{Id: 1, Title: "Task 1", Status: false},
}

func SliceTaskToJson() []byte {
	sliceJson, errorJson := json.Marshal(tasks)
	if errorJson != nil {
		fmt.Println("Error: ", errorJson)
	}
	return sliceJson
}

func FindTaskByID(id int) (*t.Task, error) {
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		if task.Id == id {
			return &task, nil
		}
	}
	return nil, errors.New("the id does not exist")
}

func main() {

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		method := r.Method
		if url == "/todos" {
			// - /todos - GET
			if method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(SliceTaskToJson())
				return
			}

			// - /todos - POST
			if method == http.MethodPost {
				body, errorBody := io.ReadAll(r.Body)
				if errorBody != nil {
					http.Error(w, errorBody.Error(), http.StatusInternalServerError)
				}
				defer r.Body.Close()
				theID := len(tasks) + 1

				var res RequestBody
				errorUnmarshal := json.Unmarshal(body, &res)
				if errorUnmarshal != nil {
					http.Error(w, errorUnmarshal.Error(), http.StatusInternalServerError)
				}

				newTask := t.NewTask(res.Title, theID)
				tasks = append(tasks, *newTask)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(SliceTaskToJson())
				return
			}
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		} else {
			http.Error(w, "The route "+url+" does not exist", http.StatusNotFound)
		}

	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		method := r.Method
		match, _ := regexp.MatchString(`/todos/\d+`, url)
		// la url debe coincidier /todos/int -> con 123s, 123/ etc intentara convertir a entero
		// - /todos/id - GET
		if match {
			if method == http.MethodGet {
				parameter := strings.TrimPrefix(r.URL.Path, "/todos/")
				id, errorAtoi := strconv.Atoi(parameter)
				if errorAtoi != nil {
					http.Error(w, errorAtoi.Error(), http.StatusInternalServerError)
					return
				}

				task, errorID := FindTaskByID(id)
				if errorID != nil {
					http.Error(w, errorID.Error(), http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(task.Json())
				return
			} else {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}

		} else {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		}

	})

	errorListen := http.ListenAndServe(":3000", nil)
	if errorListen != nil {
		log.Fatal("Listen And Serve, ", errorListen)
	}

}
