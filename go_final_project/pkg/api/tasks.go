package nextdate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dbase "github.com/avshmakov/FINAL/pkg/db"
)

type TasksResp struct {
	Tasks []*dbase.Task `json:"tasks"`
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	type errorjson struct {
		Error string `json:"error"`
	}

	tasks, err := dbase.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "dbase.Tasks(50)"})
		return
	}
	if tasks == nil {
		tasks = []*dbase.Task{} // Преобразуем nil в пустой срез
	}

	// Создаем структуру ответа с ключом "tasks"
	response := TasksResp{Tasks: tasks}

	resp, err := json.MarshalIndent(response, "", "    ")
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
	//writeJson(w, TasksResp{	Tasks: tasks,})

}

// ******************************
func getoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("inside getoneTaskHandler")

	type errorjson struct {
		Error string `json:"error"`
	}
	id := r.URL.Query().Get("id") // Теперь получаем id из query-параметра
	//fmt.Println("id=", id)

	task, err := dbase.GetTask(id) // в параметре максимальное количество записей
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "getoneTaskHandler"})
		return
	}
	resp, err := json.MarshalIndent(task, "", "    ")
	// в заголовок записываем тип контента,  данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)

}

// ********************************************************
func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	//fmt.Println("------------------putaskkHandler------------------")
	type errorjson struct {
		Error string `json:"error"`
	}
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "read body error"})

		//http.Error(w, err.Error(), http.StatusBadRequest)
		//fmt.Println("ошибка putTaskHandler тело запроса ")
		return
	}
	fmt.Println("r.body = ", buf.String())
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateAndAdjustTask(&task, w) {
		return // Ошибка уже обработана в функции
	}

	err = dbase.UpdateTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "UpdateTask"})
		return
	}
	//resp, err := json.MarshalIndent(task, "", "    ")
	// в заголовок записываем тип контента,  данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte("{}"))
}

// *******************************************
func delTaskHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("------------------delTaskHandler------------------")
	type errorjson struct {
		Error string `json:"error"`
	}

	idStr := r.URL.Query().Get("id")
	fmt.Println("idStr= ", idStr)
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "ID query error"})
		return
	}
	err := dbase.DelTask(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "DelTask error"})
		return
	}
	fmt.Println("OK")
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte("{}"))

}

// *******************************************
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("------------------doneTaskHandler------------------")
	type errorjson struct {
		Error string `json:"error"`
	}

	idStr := r.URL.Query().Get("id")
	//fmt.Println("idStr= ", idStr)
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "ID query error"})
		return
	}

	task, err := dbase.GetTask(idStr) // в параметре максимальное количество записей
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorjson{Error: "getoneTaskHandler"})
		return
	}
	//Одноразовая задача с пустым полем repeat удаляется.
	if len(task.REPEAT) == 0 {
		err := dbase.DelTask(idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorjson{Error: "DelTask error"})
			return
		}
	} else {
		var nextd string
		now := time.Now()
		nextd, err = NextDate(now, task.Date, task.REPEAT)
		err = dbase.UpdateDate(nextd, idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorjson{Error: "DelTask error"})
			return
		}

	}
	fmt.Println("OK")
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte("{}"))

}
