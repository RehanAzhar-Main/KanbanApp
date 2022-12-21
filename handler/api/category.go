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

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	ctxUserId := r.Context().Value("id")

	// check context empty
	if ctxUserId.(string) == "" {
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

	// cal category service GetCategory
	convUserid, _ := strconv.Atoi(ctxUserId.(string))
	categories, err := c.categoryService.GetCategories(r.Context(), convUserid)
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

	dataInJson, _ := json.Marshal(categories)
	// if err != nil {
	// 	panic(err)
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //200
	w.Write(dataInJson)
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	// type empty
	if category.Type == "" {
		var respSucc entity.ErrorResponse
		respSucc.Error = "invalid category request"
		dataInJson, _ := json.Marshal(respSucc)
		// if err != nil {
		// 	panic(err)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write(dataInJson)
		return
	}

	// get user_id from context
	ctxUserId, _ := strconv.Atoi(fmt.Sprintf("%s", r.Context().Value("id")))
	// if err != nil {
	// 	panic(err)
	// }

	// check context empty
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

	newCategory := entity.Category{
		Type:   category.Type,
		UserID: ctxUserId,
	}

	// call category service StoreCategory
	categoryStored, err := c.categoryService.StoreCategory(r.Context(), &newCategory)
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
		"user_id":     ctxUserId,
		"category_id": categoryStored.ID,
		"message":     "success create new category",
	})
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id from context
	ctxUserId := r.Context().Value("id")

	// get category_id
	categoryID, _ := strconv.Atoi(r.URL.Query().Get("category_id"))
	// if err != nil {
	// 	panic(err)
	// }
	err := c.categoryService.DeleteCategory(r.Context(), categoryID)
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
		"user_id":     ctxUserId,
		"category_id": categoryID,
		"message":     "success delete category",
	})
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
