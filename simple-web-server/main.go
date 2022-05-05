package main

import (
	"fmt"
  "log"
  "net/http"
)


func formHandler(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
    fmt.Fprintf(w, "Parseform() error : %v", err)
    return
  }
  fmt.Fprintf(w, "Post request successfull")
  name := r.FormValue("name")
  address := r.FormValue("address")
  fmt.Fprintf(w, "Name = %s\n", name)
  fmt.Fprintf(w, "Address = %s\n", address)

  fmt.Printf("Yes..The Form page is working\n")
  
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/hello" {
    http.Error(w, "404 Error : Page not found", http.StatusNotFound)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Method not supported", http.StatusNotFound)
    return
  }

  fmt.Printf("Yes..The Hello World page is working\n")
}




func main() {
	fmt.Println("Hello, World!")
  fileserver := http.FileServer(http.Dir("./static"))   //shorthand operator to assign and decalre variables
  http.Handle("/", fileserver)
  http.HandleFunc("/form", formHandler)
  http.HandleFunc("/hello", helloHandler)

  fmt.Printf("Server started at port 8080\n")
  if err := http.ListenAndServe(":8080", nil); err != nil{
    log.Fatal(err)
  }
}
