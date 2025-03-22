package handlers

import (
	"encoding/json"
	"huinya/db"
	"huinya/models"
	"huinya/utils"
	"net/http"
)

func UploadPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "application/json")

		var user models.User
		var post models.Post

		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, ok := utils.GetIDFromContext(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Unauthorized",
			})
			return
		}

		result := db.DB.Where("ID = ?", id).First(&user)
		if result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Not Found",
			})
			return
		}

		post.User = user

		if post.Description == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Description is required",
			})
			return
		}

		if post.Name == "" {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Name is required",
			})
			return
		}

		createPost := db.DB.Create(&post)
		if createPost.Error != nil {
			http.Error(w, `{"message" : "Create Error"}`, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Post Created",
		})

	} else {
		http.Error(w, `{"message" : "Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}
