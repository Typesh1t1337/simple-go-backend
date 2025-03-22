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

		accessToken, refreshToken, tokenErr := utils.GetJwtToken(user.ID)

		if tokenErr != nil {
			http.Error(w, tokenErr.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
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

		accessToken, refreshToken, tokenErr := utils.GetJwtToken(foundUser.ID)

		if tokenErr != nil {
			http.Error(w, tokenErr.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})

	} else {
		http.Error(w, `{"message" : "Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		id, ok := utils.GetIDFromContext(r)

		if !ok {
			http.Error(w, `{"message" : "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var user models.User

		result := db.DB.Preload("Posts").First(&user, id)
		if result.Error != nil {
			http.Error(w, `{"message" : "Email or Password is incorrect"}`, http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"email":     user.Email,
			"full_name": user.FirstName + " " + user.LastName,
			"posts":     user.Posts,
		})

	} else {
		http.Error(w, `{"message" : "Method Not Allowed"}`, http.StatusMethodNotAllowed)
	}
}
