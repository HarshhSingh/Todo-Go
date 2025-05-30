package DbHelper

import (
	"database/sql"
	"fmt"
	"main/Database"
	"net/http"
)

func GetAllTasks(body interface{}) error {

	SQL := `SELECT * FROM users`
	err := Database.Todo.Select(body, SQL)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return nil
}

func PostTaskByBody(Task string, TaskStatus string, DueDate string) error {
	SQL := fmt.Sprintf("INSERT into users(task, task_status, due_date) values($1, $2, $3)", Task, TaskStatus, DueDate)
	_, err := Database.Todo.Exec(SQL)

	return err
}

func EditTaskById(id string, Task string, TaskStatus string, DueDate string) error {
	SQL := fmt.Sprintf("UPDATE users SET task=$1, task_status=$2, due_Date=$3 where id= $4", Task, TaskStatus, DueDate, id)
	_, err := Database.Todo.Exec(SQL)
	return err
}
func IsTaskIdValid(res http.ResponseWriter, taskId string) bool {
	err := Database.Todo.Get("select * from users where id=$1", taskId)
	if err != nil {
		http.Error(res, "Failed to fetch task: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}
