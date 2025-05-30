package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"main/Database"
	"main/Database/DbHelper"
	"main/utils"
	"net/http"
)

type tasks struct {
	ID         string      `json:"id" db:"id"`
	Tasks      string      `json:"tasks" db:"task"`
	TaskStatus string      `json:"taskStatus" db:"task_status"`
	Date       string      `json:"date" db:"date"`
	DueDate    string      `json:"dueDate" db:"due_date"`
	ArchivedAt pq.NullTime `json:"archivedAt" db:"archived_at"`
}
type taskBody struct {
	Task       string `json:"task" binding:"required"`
	TaskStatus string `json:"taskStatus" binding:"required"`
	DueDate    string `json:"dueDate" binding:"required"`
}

func GetTasks(res http.ResponseWriter, req *http.Request) {
	var taskType []tasks

	err := DbHelper.GetAllTasks(&taskType)
	if err != nil {
		http.Error(res, "Failed to fetch tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("response : %v\n", taskType)
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(taskType)
	if err != nil {
		http.Error(res, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)

	fmt.Printf("Response %s \n", taskType)
	return
}

func PostTask(w http.ResponseWriter, r *http.Request) {

	var taskBody taskBody
	err := utils.DecodeJSONBody(r, &taskBody)

	fmt.Printf("payload for /task api %v \n", taskBody)
	if err != nil {
		http.Error(w, "Failed to parse the body: "+err.Error(), http.StatusBadRequest)
		return
	}
	err1 := DbHelper.PostTaskByBody(taskBody.Task, taskBody.TaskStatus, taskBody.DueDate)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task Saved"))
	fmt.Printf("Inserted /task api %v \n", taskBody)

}
func EditTask(res http.ResponseWriter, req *http.Request) {
	var taskBody taskBody
	if err := utils.DecodeJSONBody(req, &taskBody); err != nil {
		http.Error(res, "Failed to parse the body: "+err.Error(), http.StatusBadRequest)
		return
	}
	taskId := chi.URLParam(req, "taskID")
	isTaskIdValid := DbHelper.IsTaskIdValid(res, taskId)
	if !isTaskIdValid {
		http.Error(res, "Invalid task ID", http.StatusBadRequest)
		return
	}
	fmt.Printf("task Id %v \n", taskId)

	err1 := DbHelper.EditTaskById(taskId, taskBody.Task, taskBody.TaskStatus, taskBody.DueDate)
	fmt.Printf("Inserted /task api %v \n", err1)
	if err1 != nil {
		http.Error(res, "Failed to parse the body: "+err1.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Task Updated"))

}

func DeleteTask(res http.ResponseWriter, req *http.Request) {
	taskId := chi.URLParam(req, "taskID")
	if taskId == "" {
		http.Error(res, "Invalid task ID", http.StatusBadRequest)
		return
	}
	isTaskIdValid := DbHelper.IsTaskIdValid(res, taskId)
	if !isTaskIdValid {
		http.Error(res, "Invalid task ID", http.StatusBadRequest)
		return
	}
	fmt.Printf("task Id %v \n", taskId)

	_, err1 := Database.Todo.Exec(`DELETE FROM users where id= $1`, taskId)
	if err1 != nil {
		http.Error(res, "Error: "+err1.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Task Deleted"))

}
