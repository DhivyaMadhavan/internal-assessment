package main

import (
	"fmt"
	"net/http"	
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(500 * time.Millisecond)
	fmt.Fprintln(w, "Hello from Primary API")
}

func main() {
	

	http.HandleFunc("/hello", helloHandler)

	http.ListenAndServe(":8081", nil)
}