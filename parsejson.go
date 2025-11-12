package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type JsonBody struct {
	Body string `json:"body"`
}

func ParseJson(w http.ResponseWriter, r *http.Request) {
	var j = &JsonBody{}

	err := json.NewDecoder(r.Body).Decode(j)
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

	if len(j.Body) > 140 {
		SendErr(w, r, 400)
		return
	}

	str := validateJSON(j.Body)

	var resp = struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: str,
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		log.Printf("error marshalling the data")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
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
