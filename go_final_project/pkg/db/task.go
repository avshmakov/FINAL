package dbase

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	TITLE   string `json:"title"`
	COMMENT string `json:"comment"`
	REPEAT  string `json:"repeat"`
}

// ************************************************
func AddTask(task *Task) (int64, error) {

	var db *sql.DB
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		//log.Println(err)
		return 0, err
	}
	defer db.Close()
	var id int64
	// определите запрос
	res, err := db.Exec("INSERT INTO scheduler (date,title,comment,repeat) VALUES ( :p2, :p3, :p4, :p5)",
		sql.Named("p2", task.Date),
		sql.Named("p3", task.TITLE),
		sql.Named("p4", task.COMMENT),
		sql.Named("p5", task.REPEAT))

	if err != nil {
		//fmt.Println("ошибка insert ")
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		//fmt.Println("ошибка LastInsertId ")
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return id, nil
}

// ************************************************
func Tasks(limit int) ([]*Task, error) {

	var res []*Task
	var db *sql.DB

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT *from scheduler ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.TITLE, &task.COMMENT, &task.REPEAT)
		if err != nil {
			return nil, err
		}
		fmt.Println(task)
		res = append(res, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil

}

// ************************************************
func GetTask(id string) (*Task, error) {
	var db *sql.DB

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT *from scheduler WHERE id = :p1", sql.Named("p1", id))

	task := &Task{}

	err = row.Scan(&task.ID, &task.Date, &task.TITLE, &task.COMMENT, &task.REPEAT)
	if err != nil {
		return nil, err
	}
	fmt.Println(task)
	return task, nil
}

// ****************************************************************
func UpdateTask(task *Task) error {
	var db *sql.DB
	//fmt.Println("inside UpdateTask")
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	//fmt.Println("task11=", task)

	row, err := db.Exec("UPDATE scheduler SET date = :p1, title = :p2, comment = :p3, repeat = :p4 WHERE id = :p5",
		sql.Named("p1", task.Date),
		sql.Named("p2", task.TITLE),
		sql.Named("p3", task.COMMENT),
		sql.Named("p4", task.REPEAT),
		sql.Named("p5", task.ID))
	if err != nil {
		//fmt.Println("error UPDATE")
		return err
	}
	count, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil

}

// ***************************************************
func DelTask(id string) error {
	var db *sql.DB

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Println("param id for delete= ", id)

	result, err := db.Exec("DELETE from scheduler WHERE id = :p1", sql.Named("p1", id))
	if err != nil {
		fmt.Println("err=", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // Ошибка при получении количества строк
	}
	if rowsAffected == 0 {
		return fmt.Errorf("запись не найдена") // Специальная ошибка "не найдено"
	}
	//fmt.Println("delete yes")
	return nil
}

// ***************************************************

func UpdateDate(next string, id string) error {
	var db *sql.DB

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()
	//fmt.Println("param id= ", id)

	row, err := db.Exec("UPDATE scheduler SET date = :p1 WHERE id = :p2",
		sql.Named("p1", next),
		sql.Named("p2", id))
	if err != nil {
		//fmt.Println("error UPDATE new date")
		return err
	}
	count, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
