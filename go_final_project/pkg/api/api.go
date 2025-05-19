package nextdate

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("*******************")
	gnow := req.URL.Query().Get("now")
	if gnow == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("now missing"))
		return
	}
	parsTime, err := time.Parse("20060102", gnow)
	if err != nil {
		fmt.Println("Ошибка при разборе даты dstart", err)
		return
	}

	gdate := req.URL.Query().Get("date")
	if gdate == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("date missing"))
		return
	}

	grepeat := req.URL.Query().Get("repeat")
	if grepeat == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("repeat missing"))
		return
	}

	str, err := NextDate(parsTime, gdate, grepeat)
	if err != nil {
		fmt.Println("Ошибка NextDate:", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("NextDate missing"))
		return
	}
	fmt.Println("новая дата:", str)
	fmt.Println("*******************")

	res.Write([]byte(str))

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	}
}

func selectTaskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	}
}

func onetaskkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("onetaskkHandler")
	id := r.URL.Query().Get("id")
	fmt.Println("id= ", id)
	getoneTaskHandler(w, r)

}

func putttaskkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------------------putttaskkHandler------------------")

	switch r.Method {
	case http.MethodPut:
		putTaskHandler(w, r)

	}
}

func deletetaskkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------------------deletetaskkHandler------------------")
	idd := r.URL.Query().Get("id")
	fmt.Println("idd= ", idd)
	switch r.Method {
	case http.MethodDelete:
		delTaskHandler(w, r)

	}
}

func donehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------------------doneHandler------------------")
	idd := r.URL.Query().Get("id")
	fmt.Println("done id= ", idd)
	switch r.Method {
	case http.MethodPost:
		doneTaskHandler(w, r)

	}
}

func Init() *chi.Mux {
	fmt.Println("chi")

	rout := chi.NewRouter()
	rout.Handle("/*", http.FileServer(http.Dir("web")))

	rout.Get("/api/nextdate", nextDayHandler)
	rout.Post("/api/task", taskHandler)
	rout.Get("/api/tasks", selectTaskHandler)
	rout.Get("/api/task", onetaskkHandler)
	rout.Put("/api/task", putttaskkHandler)
	rout.Delete("/api/task", deletetaskkHandler)
	rout.Post("/api/task/done", donehandler)

	return rout
}

/*
func Init() {

	//http.HandleFunc("/api/nextdate", nextDayHandler)
	//http.HandleFunc("/api/task", taskHandler)
	//http.HandleFunc("/api/tasks", selectTaskHandler)
	//http.HandleFunc("/api/task/{id}", onetaskkHandler)
    r := chi.NewRouter()
	r.Get("/api/nextdate", nextDayHandler)
	r.Post("/api/task", taskHandler)
	r.Get("/api/tasks", selectTaskHandler)
	r.Get("/api/taskq/{id}", onetaskkHandler)
}
*/
