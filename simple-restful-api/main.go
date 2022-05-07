package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "github.com/gorilla/mux"
)

type Article struct {
  Title   string `json:"Title"`
  Desc    string `json:"desc"`
  Content string `json:"content"`
}

type Articles []Article

func allArticles(w http.ResponseWriter, r *http.Request){
  articles := Articles{
    Article{Title:"Test Title", Desc: "Test Description", Content:"Hey there"},
  }
  fmt.Println("Endpoint Hit: allArticles")
  json.NewEncoder(w).Encode(articles)
}

func testPostArticles(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Welcome to the Post Articles Page!")
  fmt.Println("Endpoint Hit: testPostArticles")
}

func homePage(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Welcome to the HomePage!")
  fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {

  myRouter := mux.NewRouter().StrictSlash(true)

  myRouter.HandleFunc("/", homePage)
  myRouter.HandleFunc("/articles", allArticles).Methods("GET")
  myRouter.HandleFunc("/articles", testPostArticles).Methods("POST")
  log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    handleRequests()
}