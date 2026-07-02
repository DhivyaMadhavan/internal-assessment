package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Secondary API")
}

func main() {

	http.HandleFunc("/hello", helloHandler)	

	http.ListenAndServe(":8082", nil)
}