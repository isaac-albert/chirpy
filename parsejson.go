package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/isaac-albert/chirpy/internal/database"
)

type MessageBody struct {
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type Message struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (api *apiConfig) ParseMessage(w http.ResponseWriter, r *http.Request) {
	var msg = &MessageBody{}

	err := json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		var e = struct {
			Err string `json:"error"`
		}{
			Err: err.Error(),
		}
		w.WriteHeader(500)
		data, err := json.Marshal(e)
		if err != nil {
			e.Err += err.Error()
		}
		w.Write(data)
	}

	if len(msg.Body) > 140 {
		SendErr(w, r, 400)
		return
	}

	str := validateJSON(msg.Body)
	var newMsg = database.CreateUserMessageParams{}
	newMsg.Body = str
	newMsg.UserID = msg.UserId

	msgData, err := api.dbQuery.CreateUserMessage(r.Context(), newMsg)
	if err != nil {
		log.Printf("user message: %v", err)
		return
	}

	var resp = &Message{}
	resp.Body = msgData.Body
	resp.CreatedAt = msgData.CreatedAt
	resp.UpdatedAt = msgData.UpdatedAt
	resp.Id = msgData.ID
	resp.UserId = msgData.UserID

	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error marshalling the data")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
}

func SendErr(w http.ResponseWriter, r *http.Request, status int) {
	var e = struct {
		Err string `json:"error"`
	}{
		Err: http.StatusText(status),
	}
	w.WriteHeader(status)
	data, err := json.Marshal(e)
	if err != nil {
		e.Err += err.Error()
	}
	w.Write(data)
}

func validateJSON(data string) string {
	var profane = []string{"kerfuffle", "sharbert", "fornax"}

	strParts := strings.Split(data, " ")

	strFinal := ""
	for _, value := range strParts {
		for _, word := range profane {
			if strings.Contains(strings.ToLower(value), word) {
				value = "****"
			}
		}
		strFinal += value + " "
	}

	strFinal = strings.TrimRight(strFinal, " ")

	return strFinal
}

func (afg *apiConfig) GetMessages(w http.ResponseWriter, r *http.Request) {
	userMsgs, err := afg.dbQuery.GetUsers(r.Context())
	if err != nil {
		log.Printf("getting users: %v", err)
		return
	}


	var users = make([]Message, len(userMsgs))

	for i, msg := range userMsgs {
		users[i].Body = msg.Body
		users[i].CreatedAt = msg.CreatedAt
		users[i].UpdatedAt = msg.UpdatedAt
		users[i].Id = msg.ID
		users[i].UserId = msg.UserID
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
    
    data, err := json.Marshal(users)
    if err != nil {
        http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
        return
    }
    
    w.Write(data)
}

func (afg *apiConfig) GetSingleMessage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	
	uu, err := uuid.Parse(id)
	if err != nil {
		log.Printf("id parse: %v", err)
		return
	}
	msg, err := afg.dbQuery.GetMessage(r.Context(), uu)
	if err != nil {
		w.WriteHeader(404)
		log.Printf("get message: %v", err)
		return
	}

	var m = &Message{}
	m.Body = msg.Body
	m.CreatedAt = msg.CreatedAt
	m.UpdatedAt = msg.UpdatedAt
	m.Id = msg.ID
	m.UserId = msg.UserID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	data, err := json.Marshal(m)
	if err != nil {
		log.Printf("msg marshal: %v", err)
		return
	}

	w.Write(data)
}
