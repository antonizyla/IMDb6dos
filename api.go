package main

import (
	"encoding/json"
	"html"
	"net/http"
)

func actorInfoHandler(w http.ResponseWriter, r *http.Request) {
	actorID := html.EscapeString(r.URL.Query().Get("id"))
	if actorID == "" {
		http.Error(w, "Missing required parameter 'id'", http.StatusBadRequest)
		return
	}
	actor := getActorDetails(actorID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}
