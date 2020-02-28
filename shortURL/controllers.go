package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Environment struct {
	Db Datastore
}

func (env *Environment)NewURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	obj, err := env.Db.NewURL(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(obj)
	w.WriteHeader(http.StatusOK)
}

func (env *Environment)GetHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	url := paramFromURL["short_link"]
	obj, err := env.Db.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, obj.Link, 301)
	_ = json.NewEncoder(w).Encode(obj)
	w.WriteHeader(http.StatusOK)
}

func (env *Environment)SetTtlHandler(w http.ResponseWriter, r *http.Request) {
	obj := model{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = env.Db.SetTtl(obj.Ttl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (env *Environment)GetAllHandler(w http.ResponseWriter, r *http.Request) {
	total, links, err := env.Db.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(total)
	_ = json.NewEncoder(w).Encode(links)
	w.WriteHeader(http.StatusOK)
}
