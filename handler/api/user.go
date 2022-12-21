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

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// check email or pass empty
	if user.Email == "" || user.Password == "" {

		var respErr entity.ErrorResponse
		respErr.Error = "email or password is empty"
		dataInJson, _ := json.Marshal(respErr)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write(dataInJson)
		return
	}

	newUserLog := entity.User{
		Email:    user.Email,
		Password: user.Password,
	}

	// call user service login
	logUser_id, err := u.userService.Login(r.Context(), &newUserLog)
	if err != nil {
		var respErr entity.ErrorResponse
		respErr.Error = "error internal server"
		dataInJson, _ := json.Marshal(respErr)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write(dataInJson)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(logUser_id),
	})

	// success response
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": logUser_id,
		"message": "login success",
	})

}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	// fullname or password empty
	if user.Fullname == "" || user.Email == "" || user.Password == "" {

		var respErr entity.ErrorResponse
		respErr.Error = "register data is empty"
		dataInJson, _ := json.Marshal(respErr)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write(dataInJson)
		return
	}

	// user service register
	newUser := entity.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	regUser, err := u.userService.Register(r.Context(), &newUser)
	if err != nil {
		fmt.Println("error user service register")
		var respErr entity.ErrorResponse
		respErr.Error = "error internal server"
		dataInJson, _ := json.Marshal(respErr)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //500
		w.Write(dataInJson)
		return

	}
	// fmt.Fprint(w, "jalan di api Register user")

	// success response
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": regUser.ID,
		"message": "register success",
	})

}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	// set http cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "",
		Value: "",
	})

	// success response
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "logout success",
	})
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
