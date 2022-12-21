package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	// get id from context
	getUserId := r.Context().Value("id")
	ctxUserId, _ := strconv.Atoi(getUserId.(string))
	// fmt.Println("ctxUserId", ctxUserId)
	// if err != nil {
	// 	panic(err)
	// }

	// check ctx nil
	if ctxUserId == 0 {
		var respSucc entity.ErrorResponse
		respSucc.Error = "invalid user id"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) //500
		w.Write(dataInJson)
		return
	}

	// get task_id
	getTaskID := r.URL.Query().Get("task_id")
	// ctxTaskID, _ := strconv.Atoi(getTaskID)
	// fmt.Println("ctxTaskID", ctxTaskID)
	// if err != nil {
	// 	panic(err)
	// }

	// check empty task_id
	if getTaskID == "" {
		listTask, err := t.taskService.GetTasks(r.Context(), ctxUserId)
		if err != nil {
			var respSucc entity.ErrorResponse
			respSucc.Error = "error internal server"
			dataInJson, _ := json.Marshal(respSucc)
			// if err != nil {
			// 	panic(err)
			// }
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError) //500
			w.Write(dataInJson)
			return
		}

		// give success response
		dataInJson, _ := json.Marshal(listTask)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200
		w.Write(dataInJson)
		return

	}
	if getTaskID != "" {
		ctxTaskID, _ := strconv.Atoi(getTaskID)
		taskGet, errGetTaskID := t.taskService.GetTaskByID(r.Context(), ctxTaskID)
		if errGetTaskID != nil {
			var respSucc entity.ErrorResponse
			respSucc.Error = "error internal server"
			dataInJson, _ := json.Marshal(respSucc)
			// if err != nil {
			// 	panic(err)
			// }
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError) //500
			w.Write(dataInJson)
			return
		}
		// give success response
		// dataInJson, _ := json.Marshal(taskGet)

		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200
		json.NewEncoder(w).Encode(taskGet)
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	// check empyt title, description dan category_id
	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		var respSucc entity.ErrorResponse
		respSucc.Error = "invalid task request"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write(dataInJson)
		return
	}

	// get user id from context
	getUserId := r.Context().Value("id")
	ctxUserId, _ := strconv.Atoi(getUserId.(string))
	// if err != nil {
	// 	panic(err)
	// }

	// check ctx nil
	if ctxUserId == 0 {
		var respSucc entity.ErrorResponse
		respSucc.Error = "invalid user id"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write(dataInJson)
		return
	}

	newTask := &entity.Task{
		Title:       task.Title,
		Description: task.Description,
		UserID:      ctxUserId,
		CategoryID:  task.CategoryID,
	}

	taskStored := entity.Task{}
	taskStored, err = t.taskService.StoreTask(r.Context(), newTask)
	if err != nil {
		var respSucc entity.ErrorResponse
		respSucc.Error = "error internal server"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write(dataInJson)
		return
	}

	// success response
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": ctxUserId,
		"task_id": taskStored.ID,
		"message": "success create new task",
	})
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	ctxUserId, _ := strconv.Atoi(fmt.Sprintf("%s", r.Context().Value("id")))
	// if err != nil {
	// 	panic(err)
	// }

	// get task_id
	ctxTaskID, _ := strconv.Atoi(r.URL.Query().Get("task_id"))
	// if err != nil {
	// 	panic(err)
	// }

	err := t.taskService.DeleteTask(r.Context(), ctxTaskID)
	if err != nil {
		var respSucc entity.ErrorResponse
		respSucc.Error = "error internal server"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write(dataInJson)
		return
	}

	// success response
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": ctxUserId,
		"task_id": ctxTaskID,
		"message": "success delete task",
	})

}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// get user id from context
	getUserId := r.Context().Value("id")

	// if err != nil {
	// 	panic(err)
	// }

	// check ctx user_id nil
	if getUserId.(string) == "" {
		var respSucc entity.ErrorResponse
		respSucc.Error = "invalid user id"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write(dataInJson)
		return
	}
	ctxUserId, _ := strconv.Atoi(getUserId.(string))
	updtTask := entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      ctxUserId,
	}

	taskUpdated, err := t.taskService.UpdateTask(r.Context(), &updtTask)
	if err != nil {
		var respSucc entity.ErrorResponse
		respSucc.Error = "error internal server"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write(dataInJson)
		return
	}

	// success response
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": ctxUserId,
		"task_id": taskUpdated.ID,
		"message": "success update task",
	})
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
