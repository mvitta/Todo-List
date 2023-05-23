package main

import (
	"encoding/json"
	"io"
	"log"
	f "main/functions"
	t "main/task"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type RequestBody struct {
	Title string `json:"title"`
}

type RequestBodyUpdate struct {
	Status bool `json:"status"`
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
				w.Write(f.SliceTasksToJson())
				return
			}
			// hacer una funcion que retorne un error y que contenga codigo repetido
			// - /todos - POST
			if method == http.MethodPost {
				body, errorBody := io.ReadAll(r.Body)
				if errorBody != nil {
					http.Error(w, errorBody.Error(), http.StatusInternalServerError)
					return
				}
				defer r.Body.Close()
				theID := len(f.Tasks) + 1

				var req RequestBody
				errorUnmarshal := json.Unmarshal(body, &req)
				if errorUnmarshal != nil {
					http.Error(w, errorUnmarshal.Error(), http.StatusInternalServerError)
					return
				}

				newTask := t.NewTask(req.Title, theID)
				f.Tasks = append(f.Tasks, *newTask)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(f.SliceTasksToJson())
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
		if match {
			parameter := strings.TrimPrefix(url, "/todos/")
			id, errorAtoi := strconv.Atoi(parameter)
			if errorAtoi != nil {
				http.Error(w, errorAtoi.Error(), http.StatusInternalServerError)
				return
			}

			// - /todos/id - GET
			if method == http.MethodGet {

				task, errorID := f.FindTaskByID(id)
				if errorID != nil {
					http.Error(w, errorID.Error(), http.StatusNotFound)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(task.Json())

				// - /todos/id - PUT
			} else if method == http.MethodPut {

				// get Body
				body, errorBody := io.ReadAll(r.Body)
				if errorBody != nil {
					http.Error(w, errorBody.Error(), http.StatusInternalServerError)
					return
				}

				defer r.Body.Close()
				var newStatus RequestBodyUpdate

				if errUnmarshal := json.Unmarshal(body, &newStatus); errUnmarshal != nil {
					http.Error(w, errUnmarshal.Error(), http.StatusInternalServerError)
					return
				}

				if errUpdate := f.FindAndUpdateStatus(id, newStatus.Status); errUpdate != nil {
					http.Error(w, errUpdate.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)

			} else {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}

		} else {
			//no muestra el string -> "404 Not Found, do you mean /todos/id?"
			http.Error(w, "404 Not Found, do you mean /todos/id?", http.StatusNotFound)
		}

	})

	if errorListen := http.ListenAndServe(":3000", nil); errorListen != nil {
		log.Fatal("Listen And Serve, ", errorListen)
	}

}
