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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodGet {
			routes := map[string]string{
				"/-GET":              "https://bc90-186-112-196-136.ngrok-free.app",
				"/todos-GET":         "https://bc90-186-112-196-136.ngrok-free.app/todos",
				"/todos-POST":        "https://bc90-186-112-196-136.ngrok-free.app/todos",
				"/todos/<id>-GET":    "https://bc90-186-112-196-136.ngrok-free.app/todos/<id>",
				"/todos/<id>-PUT":    "https://bc90-186-112-196-136.ngrok-free.app/todos/<id>",
				"/todos/<id>-DELETE": "https://bc90-186-112-196-136.ngrok-free.app/todos/<id>",
				"how does it work":   "https://docs.google.com/document/d/1x4qcJdOrMg3ET55nYUcPBVFEdqY8SgmqAF0KBeD5f4c/edit",
			}
			if json, errorJson := json.Marshal(routes); errorJson != nil {
				http.Error(w, errorJson.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(json)
			}
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	})

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
			switch method {
			// - /todos/id - GET
			case http.MethodGet:
				task, errorID := f.FindTaskByID(id)
				if errorID != nil {
					http.Error(w, errorID.Error(), http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(task.Json())

			// - /todos/id - PUT
			case http.MethodPut:
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
				w.WriteHeader(http.StatusNoContent)
			// - /todos/id - DELETE
			case http.MethodDelete:
				if errorDelete := f.DeleteTask(id); errorDelete != nil {
					http.Error(w, errorDelete.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.Error(w, "404 Not Found, do you mean /todos/id -> id must be an integer?", http.StatusNotFound)
		}
	})

	if errorListen := http.ListenAndServe(":3000", nil); errorListen != nil {
		log.Fatal("Listen And Serve, ", errorListen)
	}

}
