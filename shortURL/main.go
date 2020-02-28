package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {
	config := ReadConfig()
	fmt.Println(config)
	fmt.Println("Connecting to database server...")

	db, err := NewDB(config)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Success connection")

	if err = db.Ping(); err != nil {
		fmt.Println(err)
	}
	env := &Environment{Db: db}

	r := mux.NewRouter()
	r.Use(SetJSONHeader)
	r.HandleFunc("/new_short_url", env.NewURLHandler).Methods("GET")
	r.HandleFunc("/{short_link}", env.GetHandler).Methods("GET")
	r.HandleFunc("/admin/set_ttl", env.SetTtlHandler).Methods("POST")
	r.HandleFunc("/admin/get_all", env.GetAllHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func SetJSONHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
