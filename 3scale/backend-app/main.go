package main

import (
	"fmt"
	"net/http"
	"os"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func listenAndServe(port string) {
	fmt.Printf("serving on %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func main() {
	http.HandleFunc("/", okHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	listenAndServe(port)
}