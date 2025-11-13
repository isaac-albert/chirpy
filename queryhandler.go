package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type email struct {
	E string `json:"email"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (afg *apiConfig) apiQueryHandler(w http.ResponseWriter, r *http.Request) {
	var e = &email{}
	err := json.NewDecoder(r.Body).Decode(e)
	if err != nil {
		log.Printf("deocode email: %v", err)
		return
	}
	userObj, err := afg.dbQuery.CreateUser(r.Context(), e.E)
	if err != nil {
		log.Printf("create user: %v", err)
		return
	}

	var user = &User{}

	user.Id = userObj.ID
	user.CreatedAt = userObj.CreatedAt
	user.UpdatedAt = userObj.UpdatedAt
	user.Email = userObj.Email

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("user marshalling: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(userJson)
}

