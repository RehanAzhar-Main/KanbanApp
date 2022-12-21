package web

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/entity"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
)

type DashboardWeb interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
}

type dashboardWeb struct {
	categoryClient client.CategoryClient
	embed          embed.FS
}

func NewDashboardWeb(catClient client.CategoryClient, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{catClient, embed}
}

func (d *dashboardWeb) Dashboard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	categories, err := d.categoryClient.GetCategories(userId.(string))
	if err != nil {
		log.Println("error get cat: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dataTemplate = map[string]interface{}{
		"categories": categories,
	}

	var getIndexByCategoryId = func(catId int) int {
		for i := 0; i < len(categories); i++ {
			if categories[i].ID == catId {
				return i
			}
		}

		return -1
	}

	var funcMap = template.FuncMap{
		"categoryInc": func(categoryId int) int {
			idx := getIndexByCategoryId(categoryId)

			if idx == len(categories)-1 {
				return categoryId
			} else {
				return categories[idx+1].ID
			}
		},
		"categoryDec": func(categoryId int) int {
			idx := getIndexByCategoryId(categoryId)

			if idx == 0 {
				return categoryId
			} else {
				return categories[idx-1].ID
			}
		},
	}

	// ignore this
	// _ = dataTemplate
	_ = funcMap
	//
	fmt.Println("dataTemplate", dataTemplate)

	//parseFS
	mainTemplate := path.Join("views", "main", "dashboard.html")
	header := path.Join("views", "general", "header.html")

	tmpl, err := template.ParseFS(d.embed, mainTemplate, header)

	err = tmpl.Execute(w, dataTemplate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.ErrorResponse{Error: err.Error()})
	}
}
