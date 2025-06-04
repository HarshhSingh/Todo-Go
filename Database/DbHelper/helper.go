package DbHelper

import (
	"database/sql"
	"fmt"
	"main/Database"
	"main/utils"
)

func GetTasksById(body interface{}, id string, status string, search string) error {
	args := []interface{}{
		id,
		status,
		search,
	}
	SQL := `SELECT id,task,task_status, date, due_date 
             FROM tasks 
             where user_id=$1 
                 and ($2 = '' or task_status = $2)
                 and ($3 = '' or task ilike '%' || $3 || '%')
                 and archived_at IS NULL
             ORDER BY date DESC`
	err := Database.Todo.Select(body, SQL, args...)
	fmt.Printf("GetTasksById err: %v\n", err)
	fmt.Printf("GetTasksById response: %v\n", body)

	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return nil
	}
	return nil
}

func PostTaskByBody(task string, taskStatus string, dueDate string, userId string) error {
	args := []interface{}{
		task,
		taskStatus,
		dueDate,
		userId,
	}
	SQL := "INSERT into tasks(task, task_status, due_date,user_id) values($1, $2, $3,$4)"
	_, err := Database.Todo.Exec(SQL, args...)
	return err
}

func EditTaskById(id string, userId string, task string, taskStatus string, dueDate string) error {
	args := []interface{}{
		task,
		taskStatus,
		dueDate,
		id, userId,
	}
	SQL := "UPDATE tasks SET task=$1, task_status=$2, due_Date=$3 where id= $4 and user_id=$5"
	_, err := Database.Todo.Exec(SQL, args...)
	return err
}
func IsTaskIdValid(taskId string) (bool, error) {
	var isValid bool
	SQL := `select count(*)>0 from tasks where id=$1`
	err := Database.Todo.Get(&isValid, SQL, taskId)
	return isValid, err
}
func IsTaskIdValidForUser(taskId string, userId string) (bool, error) {
	var isValid bool

	SQL := "select count(*)>0 from tasks where id=$1 and user_id=$2"
	err := Database.Todo.Get(&isValid, SQL, taskId, userId)
	return isValid, err
}
func IsUserIDValid(userId string) (bool, error) {
	var isValid bool
	SQL := `SELECT count(id)>0 FROM users WHERE id=$1 AND archived_at IS NULL`
	err := Database.Todo.Get(&isValid, SQL, userId)
	return isValid, err
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
