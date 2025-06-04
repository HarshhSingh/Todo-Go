package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"main/Database"
	"main/Database/DbHelper"
	"main/constants"
	"main/utils"
	"net/http"
	"strings"
	"time"
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
type userBody struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type loginBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(res http.ResponseWriter, req *http.Request) {
	var payload userBody
	err := utils.DecodeJSONBody(req, &payload)
	fmt.Printf("payload for /register api %v \n", payload)
	if err != nil {
		fmt.Printf("RegisterUser: error while registering user: %v.", err)
		http.Error(res, "Failed to parse the body: "+err.Error(), http.StatusBadRequest)
		return
	}
	exists, existsErr := DbHelper.IfEmailExists(payload.Email)

	if existsErr != nil {
		http.Error(res, "failed to check user existence:"+existsErr.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(res, "user already exists", http.StatusBadRequest)
		return
	}
	hashedPassword, err := utils.HashPassword(payload.Password)

	if err != nil {
		http.Error(res, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err1 := DbHelper.RegisterUser(payload.Name, payload.Email, hashedPassword)
	if err1 != nil {
		http.Error(res, "Failed to register: "+err1.Error(), http.StatusInternalServerError)
		return
	}
	//res.WriteHeader(http.StatusOK)

	utils.RespondJSON(res, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "User registered",
	})

}
func LoginUser(res http.ResponseWriter, req *http.Request) {
	var payload loginBody
	err := utils.DecodeJSONBody(req, &payload)
	fmt.Printf("payload for /login api %v \n", payload)
	if err != nil {
		http.Error(res, "Failed to parse the body: "+err.Error(), http.StatusBadRequest)
		return
	}
	userId, err1 := DbHelper.LoginUser(payload.Email, payload.Password)
	if err1 != nil {
		http.Error(res, "Failed to login: "+err1.Error(), http.StatusInternalServerError)
		return
	}
	sessionToken, tokenErr := utils.CreateToken(userId + time.Now().String())
	if tokenErr != nil {
		http.Error(res, "Failed to create token: "+tokenErr.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("sessionToken: %v\n userId: %v\n", sessionToken, userId)
	utils.RespondJSON(res, http.StatusCreated, struct {
		Token   string `json:"token"`
		Message string `json:"message"`
	}{
		Message: "Login successful",
		Token:   sessionToken,
	})
	res.WriteHeader(http.StatusOK)
}
func IsValidStatus(status string) bool {
	fmt.Printf("status: %v\n", status)

	for _, validStatus := range constants.AllTaskStatuses {
		if strings.EqualFold(status, validStatus) { // case-insensitive match
			return true
		}
	}
	return false

}
func AllTasksById(res http.ResponseWriter, req *http.Request) {
	var taskType []tasks
	userId := chi.URLParam(req, "userId")
	queryParamsStatus := req.URL.Query().Get("status")
	queryParamsSearch := req.URL.Query().Get("search")

	if queryParamsStatus != "" {
		statusBool := IsValidStatus(queryParamsStatus)
		if !statusBool {
			http.Error(res, "Invalid status", http.StatusBadRequest)
			return
		}
	}
	res.Header().Set("Content-Type", "application/json")
	fmt.Printf("userId: %v\n", userId)
	validIdBool, IsUserIDValidErr := DbHelper.IsUserIDValid(userId)
	if !validIdBool {
		fmt.Printf("Invalid user id: %v\n", userId)
		http.Error(res, "Invalid user id: "+userId, http.StatusBadRequest)
		return
	}
	fmt.Printf("bool: %v\n", validIdBool)
	if IsUserIDValidErr != nil {
		fmt.Printf("IsUserIDValidErr: %v\n", IsUserIDValidErr)
		http.Error(res, "Error: ", http.StatusBadRequest)
		return
	}
	err := DbHelper.GetTasksById(&taskType, userId, queryParamsStatus, queryParamsSearch)
	if err != nil {
		http.Error(res, "Failed to fetch tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("response : %v\n", taskType)
	res.Header().Set("Content-Type", "application/json")

	if taskType == nil || len(taskType) == 0 {
		utils.RespondJSON(res, http.StatusOK, struct {
			Tasks []tasks `json:"tasks"`
		}{
			Tasks: []tasks{},
		})
		return
	}
	utils.RespondJSON(res, http.StatusOK, struct {
		Tasks []tasks `json:"tasks"`
	}{Tasks: taskType})
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
	userId := chi.URLParam(r, "userId")
	fmt.Printf("userId: %v\n", userId)
	validIdBool, IsUserIDValidErr := DbHelper.IsUserIDValid(userId)
	if !validIdBool {
		fmt.Printf("Invalid user id: %v\n", userId)
		http.Error(w, "Invalid user id: "+userId, http.StatusBadRequest)
		return
	}
	fmt.Printf("bool: %v\n", validIdBool)
	if IsUserIDValidErr != nil {
		fmt.Printf("IsUserIDValidErr: %v\n", IsUserIDValidErr)
		http.Error(w, "Error: ", http.StatusBadRequest)
		return
	}
	err1 := DbHelper.PostTaskByBody(taskBody.Task, taskBody.TaskStatus, taskBody.DueDate, userId)
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
	isTaskIdValid, isTaskIdValidErr := DbHelper.IsTaskIdValid(taskId)
	if isTaskIdValidErr != nil {
		fmt.Printf("isTaskIdValidErr: %v\n", isTaskIdValidErr)
		http.Error(res, "Error: "+isTaskIdValidErr.Error(), http.StatusBadRequest)
		return
	}

	if !isTaskIdValid {
		fmt.Printf("Invalid task id: %v\n", taskId)
		http.Error(res, "Invalid task ID", http.StatusBadRequest)
		return
	}
	fmt.Printf("task Id %v \n", taskId)

	userId := chi.URLParam(req, "userId")
	fmt.Printf("userId: %v\n", userId)
	validIdBool, IsUserIDValidErr := DbHelper.IsUserIDValid(userId)
	if !validIdBool {
		fmt.Printf("Invalid user id: %v\n", userId)
		http.Error(res, "Invalid user id: "+userId, http.StatusBadRequest)
		return
	}
	fmt.Printf("bool: %v\n", validIdBool)
	if IsUserIDValidErr != nil {
		fmt.Printf("IsUserIDValidErr: %v\n", IsUserIDValidErr)
		http.Error(res, "Error: ", http.StatusBadRequest)
		return
	}
	isTaskIdValidForUser, isTaskIdValidForUserErr := DbHelper.IsTaskIdValidForUser(taskId, userId)
	if !isTaskIdValidForUser {
		http.Error(res, "Invalid task ID", http.StatusUnauthorized)
		return
	}
	if isTaskIdValidForUserErr != nil {
		http.Error(res, "Error: "+isTaskIdValidForUserErr.Error(), http.StatusBadRequest)
		return
	}
	err1 := DbHelper.EditTaskById(taskId, userId, taskBody.Task, taskBody.TaskStatus, taskBody.DueDate)
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
	isTaskIdValid, isTaskIdValidErr := DbHelper.IsTaskIdValid(taskId)
	if isTaskIdValidErr != nil {
		fmt.Printf("isTaskIdValidErr: %v\n", isTaskIdValidErr)
		http.Error(res, "Error: "+isTaskIdValidErr.Error(), http.StatusBadRequest)
		return
	}

	if !isTaskIdValid {
		fmt.Printf("Invalid task id: %v\n", taskId)
		http.Error(res, "Invalid task ID", http.StatusBadRequest)
		return
	}
	userId := chi.URLParam(req, "userId")
	fmt.Printf("userId: %v\n", userId)
	validIdBool, IsUserIDValidErr := DbHelper.IsUserIDValid(userId)
	if !validIdBool {
		fmt.Printf("Invalid user id: %v\n", userId)
		http.Error(res, "Invalid user id: "+userId, http.StatusBadRequest)
		return
	}
	fmt.Printf("bool: %v\n", validIdBool)
	if IsUserIDValidErr != nil {
		fmt.Printf("IsUserIDValidErr: %v\n", IsUserIDValidErr)
		http.Error(res, "Error: ", http.StatusBadRequest)
		return
	}
	isTaskIdValidForUser, isTaskIdValidForUserErr := DbHelper.IsTaskIdValidForUser(taskId, userId)
	if !isTaskIdValidForUser {
		http.Error(res, "Invalid task ID", http.StatusUnauthorized)
		return
	}
	if isTaskIdValidForUserErr != nil {
		http.Error(res, "Error: "+isTaskIdValidForUserErr.Error(), http.StatusBadRequest)
		return
	}
	_, err1 := Database.Todo.Exec(`DELETE FROM tasks where id= $1 and user_id=$2`, taskId, userId)
	if err1 != nil {
		http.Error(res, "Error: "+err1.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Task Deleted"))

}
