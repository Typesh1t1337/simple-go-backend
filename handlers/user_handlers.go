package handlers

import (
	"encoding/json"
	"huinya/db"
	"huinya/models"
	"huinya/utils"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, `{"message" : "JSON Error"}`, http.StatusBadRequest)
			return
		}

		if user.Email == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" {
			http.Error(w, `{"message" : "Password, Email, FirstName or Lastname is required"}`, http.StatusBadRequest)
			return
		}

		hashedPassword, errHash := utils.HashPassword(user.Password)

		if errHash != nil {
			http.Error(w, `{"message" : "Hash Error"}`, http.StatusBadRequest)
			return
		}

		user.Password = hashedPassword

		result := db.DB.Create(&user)

		if result.Error != nil {
			http.Error(w, `{"message" : "Create Error"}`, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User created successfully",
		})

	} else {
		http.Error(w, `{"message" : "This method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var user models.User
		var foundUser models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, `{"message" : "JSON Error"}`, http.StatusBadRequest)
			return
		}

		if user.Email == "" || user.Password == "" {
			http.Error(w, `{"message" : "Email and Password are required"}`, http.StatusBadRequest)
			return
		}

		result := db.DB.Where("email = ?", user.Email).First(&foundUser)

		if result.Error != nil {
			http.Error(w, `{"message" : "email or password is incorrect"}`, http.StatusNotFound)
			return
		}

		hashErr := utils.CheckPasswordHash(user.Password, foundUser.Password)
		if hashErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Email or Password is incorrect",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User logged successfully",
		})

	} else {
		http.Error(w, `{"message" : "Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}
