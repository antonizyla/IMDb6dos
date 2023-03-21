package main

import (
	"encoding/json"
	"html"
	"net/http"
)

func actorInfoHandler(w http.ResponseWriter, r *http.Request) {
	actorID := html.EscapeString(r.URL.Query().Get("nconst"))
	if actorID == "" {
		http.Error(w, "Missing required parameter 'nconst'", http.StatusBadRequest)
		return
	}
	actor, err := getActorDetails(actorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}

func titlesInfoHandler(w http.ResponseWriter, r *http.Request) {
    titleID := html.EscapeString(r.URL.Query().Get("tconst"))
    if titleID == "" {
        http.Error(w, "Missing required parameter 'tconst'", http.StatusBadRequest)
        return
    }
    title, err := getTitleDetails(titleID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(title)
}
