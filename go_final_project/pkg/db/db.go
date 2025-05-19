package dbase

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""

);
CREATE INDEX date_index ON scheduler (date)
`

// Init(dbFile string) error
func Init(dbFile string) error {
	var DB *sql.DB
	var err error
	_, err = os.Stat(dbFile)

	if err != nil {
		fmt.Println("Database is absent, creating scheduler.db")
		file, err := os.Create(dbFile)
		if err != nil {
			return err
		}
		defer file.Close()
		fmt.Println("scheduler.db was created")
		// Открываем соединение с БД
		DB, err := sql.Open("sqlite", dbFile)
		if err != nil {
			fmt.Println("failed to open DB")
			return err
		}
		defer DB.Close()

		_, err = DB.Exec(schema)
		if err != nil {
			fmt.Println("failed to create table")
			return err
		}
		fmt.Println("Table scheduler was created!")
	} else {
		fmt.Println("Database ok")
		DB, err = sql.Open("sqlite", dbFile)
		if err != nil {
			return err
		}
		defer DB.Close()
	}
	return nil
}
