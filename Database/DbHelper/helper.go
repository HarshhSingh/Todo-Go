package DbHelper

import (
	"database/sql"
	"fmt"
	"main/Database"
	"main/utils"
	"net/http"
)

func GetAllTasks(body interface{}) error {

	SQL := `SELECT * FROM tasks`
	err := Database.Todo.Select(body, SQL)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return nil
}

func PostTaskByBody(task string, taskStatus string, dueDate string) error {
	args := []interface{}{
		task,
		taskStatus,
		dueDate,
	}
	SQL := "INSERT into tasks(task, task_status, due_date) values($1, $2, $3)"
	_, err := Database.Todo.Exec(SQL, args...)
	return err
}

func EditTaskById(id string, task string, taskStatus string, dueDate string) error {
	args := []interface{}{
		task,
		taskStatus,
		dueDate,
	}
	SQL := "UPDATE tasks SET task=$1, task_status=$2, due_Date=$3 where id= $4"
	_, err := Database.Todo.Exec(SQL, args...)
	return err
}
func IsTaskIdValid(res http.ResponseWriter, taskId string) bool {
	err := Database.Todo.Get("select * from tasks where id=$1", taskId)
	if err != nil {
		http.Error(res, "Failed to fetch task: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}
func IfEmailExists(email string) (bool, error) {
	SQL := `SELECT id FROM users WHERE email = TRIM(LOWER($1)) AND archived_at IS NULL`
	var id string
	err := Database.Todo.Get(&id, SQL, email)
	fmt.Printf("id: %v\n", id)
	fmt.Printf("err: %v\n", err)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func RegisterUser(name string, email string, password string) error {
	SQL := "INSERT into users(name, email, password) values($1, $2, $3)"
	args := []interface{}{
		name,
		email,
		password,
	}
	_, err := Database.Todo.Exec(SQL, args...)
	return err
}
func LoginUser(email string, password string) (string, error) {
	SQL := "SELECT id,password FROM users WHERE email=$1 AND archived_at IS NULL"
	var hashedPassword string
	var userId string

	err := Database.Todo.QueryRowx(SQL, email).Scan(&userId, &hashedPassword)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	fmt.Printf("hashedPassword: %v\n realPassword: %v\n", hashedPassword, password)
	if err == sql.ErrNoRows {
		return "", nil
	}
	checkPasswordErr := utils.CheckPassword(password, hashedPassword)
	fmt.Printf("checkPassword: %v\n", checkPasswordErr)
	if checkPasswordErr != nil {
		return "", fmt.Errorf("invalid password")
	}

	return userId, nil
}
