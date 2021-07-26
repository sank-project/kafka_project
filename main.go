package main

import (
	"encoding/json"
	"fmt"
	"kafka/model"
	"kafka/producer"
	restfunction "kafka/rest_function"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var posts []restfunction.Post

func main() {
	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	router := mux.NewRouter()

	posts = append(posts, restfunction.Post{ID: "1", Title: "My First Post", Body: "This is the content of my fist post"})

	restfunction.AmbilSatu(posts)

	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")

	config := model.GetConfiguration()
	fmt.Println("Server is run, listening at port 8000")
	http.ListenAndServe(config.HttpPort, router)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range posts {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post restfunction.Post
	_ = json.NewDecoder(r.Body).Decode(&post)

	post.ID = strconv.Itoa(rand.Intn(1000000))
	posts = append(posts, post)
	print("Data Successfully Inserted!!\n\n")

	msg := fmt.Sprintf("|ID:%s|BODY:%s|TITLE:%s", post.ID, post.Body, post.Title)

	err := producer.SentToKafka(msg)
	if err != nil {
		model.GetInternalError(err, w)
		return
	}

	json.NewEncoder(w).Encode(&post)

}

func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range posts {
		if item.ID == params["id"] {
			posts = append(posts[:index], posts[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(posts)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range posts {
		if item.ID == params["id"] {
			posts = append(posts[:index], posts[index+1:]...)

			var post restfunction.Post
			_ = json.NewDecoder(r.Body).Decode(&post)

			post.ID = params["id"]
			posts = append(posts, post)
			json.NewEncoder(w).Encode(&post)
			return
		}
	}
	json.NewEncoder(w).Encode(posts)
}
