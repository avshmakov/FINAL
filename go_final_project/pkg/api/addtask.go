package nextdate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dbase "go1f/pkg/db" //"github.com/avshmakov/FINAL/pkg/db"
)

var task dbase.Task

func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	type errorjson struct {
		Error string `json:"error"`
	}

	var buf bytes.Buffer
	// читаем тело запроса

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		//fmt.Println("ошибка тело запроса ")
		return
	}
	// десериализуем JSON
	fmt.Println("*****************************************************************")
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, "ошибка десериализации JSON", http.StatusBadRequest)
		//fmt.Println("ошибка десериализации JSON = ")
		return
	}

	if !validateAndAdjustTask(&task, w) {
		return // Ошибка уже обработана в функции
	}

	id, err := dbase.AddTask(&task)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		//json.NewEncoder(w).Encode(errorjson{Error: "ошибка AddTask",})
		if err = json.NewEncoder(w).Encode(errorjson{Error: "ошибка AddTask"}); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	if err = json.NewEncoder(w).Encode(map[string]interface{}{"id": id}); err != nil {
		http.Error(w, "error response", http.StatusInternalServerError)
	}
	return
}

func validateAndAdjustTask(task *dbase.Task, w http.ResponseWriter) bool {
	type errorjson struct {
		Error string `json:"error"`
	}

	fmt.Println("task.TITLE = ", task.Title)
	fmt.Println("task.COMMENT = ", task.Comment)
	fmt.Println("task.Date = ", task.Date)
	fmt.Println("task.Repeat = ", task.Repeat)

	// Проверяем пустой или нет заголовок title
	if len(task.Title) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorjson{Error: "Не указан заголовок задачи"})
		return false
	}

	// Проверить на корректность полученное значение task.Date
	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format(DateFormat)
	} else {
		t, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorjson{Error: "дата представлена в формате, отличном от 20060102"})
			return false
		}

		var nextd string
		if len(task.Repeat) > 0 {
			fmt.Println("input NextDate ", task.Date)
			nextd, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errorjson{Error: "ошибка NextDate"})
				return false
			}
			fmt.Println("output NextDate ", nextd)
		}

		if afterNow(now, t) {
			fmt.Println("if now.After(t)")
			if len(task.Repeat) == 0 {
				task.Date = now.Format(DateFormat)
				fmt.Println("if len(task.REPEAT) == 0", task.Date)
			} else {
				task.Date = nextd
				fmt.Println("task.Date = nextd", nextd)
			}
		}
	}
	return true
}
