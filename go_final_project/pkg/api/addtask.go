package nextdate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dbase "github.com/avshmakov/FINAL/pkg/db"
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

	type answerjson struct {
		id string `json:"id"`
	}

	var buf bytes.Buffer
	// читаем тело запроса

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("ошибка тело запроса ")
		return
	}
	//fmt.Println("*****************************************************************")
	//fmt.Println("r.body = ", buf.String())

	// десериализуем JSON
	fmt.Println("*****************************************************************")
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, "ошибка десериализации JSON", http.StatusBadRequest)
		fmt.Println("ошибка десериализации JSON = ")
		return
	}

	if !validateAndAdjustTask(&task, w) {
		return // Ошибка уже обработана в функции
	}

	//+++++
	/*
		fmt.Println("task.TITLE = ", task.TITLE)
		fmt.Println("task.COMMENT = ", task.COMMENT)
		fmt.Println("task.Date = ", task.Date)
		fmt.Println("task.Repeat = ", task.REPEAT)

		// Проверяем пустой или нет  заголовок title
		if len(task.TITLE) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorjson{Error: "Не указан заголовок задачи"})
			return
		}

		// Проверить на корректность полученное значение task.Date
		now := time.Now()
		//var t,nextd time.Ticker
		//если task.Date пустая строка, то присваиваем ему текущее время now.Format("20060102");
		if len(task.Date) == 0 {
			task.Date = now.Format("20060102")
		} else {
			//проверяем, что в task.Date указана корректная дата t, err := time.Parse("20060102", task.Date). t нам ещё пригодится
			t, err := time.Parse("20060102", task.Date)
			if err != nil {

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errorjson{Error: "дата представлена в формате, отличном от 20060102"})
				return
			}
			//fmt.Println("Полет нормальный 1")
			var nextd string // Объявляем переменную заранее
			//если определён task.Repeat, то проверяем корректность правила и заодно получаем следующую дату next, err = NextDate(now, task.Date, task.Repeat);
			if len(task.REPEAT) > 0 {
				fmt.Println("input NextDate ", task.Date)
				nextd, err = NextDate(now, task.Date, task.REPEAT)
				if err != nil {

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(errorjson{Error: "ошибка  NextDate"})
					return

				}
				fmt.Println("output NextDate ", nextd)
			}

			//fmt.Println("++++++")
			//fmt.Println("now = ", now)
			//fmt.Println("t = ", t)
			//fmt.Println("++++++")

			//if now.After(t) {
			if afterNow(now, t) {
				fmt.Println("if now.After(t)")
				if len(task.REPEAT) == 0 {
					// если правила повторения нет, то берём сегодняшнее число
					task.Date = now.Format("20060102")
					fmt.Println("if len(task.REPEAT) == 0", task.Date)
				} else {
					// в противном случае, берём вычисленную ранее следующую дату
					task.Date = nextd
					fmt.Println("task.Date = nextd", nextd)
				}
			}

		}
	*/
	//+++++
	//fmt.Println("Полет нормальный 3")

	id, err1 := dbase.AddTask(&task)
	if err1 != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorjson{
			Error: "ошибка AddTask",
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})

}

func validateAndAdjustTask(task *dbase.Task, w http.ResponseWriter) bool {
	type errorjson struct {
		Error string `json:"error"`
	}

	fmt.Println("task.TITLE = ", task.TITLE)
	fmt.Println("task.COMMENT = ", task.COMMENT)
	fmt.Println("task.Date = ", task.Date)
	fmt.Println("task.Repeat = ", task.REPEAT)

	// Проверяем пустой или нет заголовок title
	if len(task.TITLE) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorjson{Error: "Не указан заголовок задачи"})
		return false
	}

	// Проверить на корректность полученное значение task.Date
	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format("20060102")
	} else {
		t, err := time.Parse("20060102", task.Date)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorjson{Error: "дата представлена в формате, отличном от 20060102"})
			return false
		}

		var nextd string
		if len(task.REPEAT) > 0 {
			fmt.Println("input NextDate ", task.Date)
			nextd, err = NextDate(now, task.Date, task.REPEAT)
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
			if len(task.REPEAT) == 0 {
				task.Date = now.Format("20060102")
				fmt.Println("if len(task.REPEAT) == 0", task.Date)
			} else {
				task.Date = nextd
				fmt.Println("task.Date = nextd", nextd)
			}
		}
	}
	return true
}
