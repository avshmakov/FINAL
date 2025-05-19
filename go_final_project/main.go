package main

import (
	"fmt"
	//nextdate "go1f/pkg/api"
	dbase "go1f/pkg/db"
	"go1f/pkg/server"
	//"time"
)

//var Install bool

func main() {
	var err error
	dbFile := "scheduler.db"

	err = dbase.Init(dbFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Запускаем сервер")
	err = server.Run() //http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}
