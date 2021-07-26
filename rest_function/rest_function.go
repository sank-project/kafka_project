package restfunction

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var Posts []Post

func AmbilSatu(Posts []Post) {
	Router := mux.NewRouter()
	Router.HandleFunc("/posts", GetPosts).Methods("GET")
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Posts)
}
