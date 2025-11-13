package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.Platfrom != "dev" {
		w.WriteHeader(403)
		return
	}
	
	val, err := cfg.dbQuery.DeleteTable(req.Context())
	if err != nil {
		w.WriteHeader(500)
		log.Printf("delete users: %v", err)
		return	
	}
	w.WriteHeader(http.StatusOK)
	data := []byte(fmt.Sprintf("rows affected: %v", val))
	w.Write(data)
}
