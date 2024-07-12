package main

import (
	"fmt"
	"net/http"

	"github.com/PCPedroso/pcp-pcp-word-gen/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/add", handlers.AddHandler)
	mux.HandleFunc("/process", handlers.ProcessHandler)
	mux.HandleFunc("/clear", handlers.ClearHandler)

	mux.HandleFunc("/gabaritojson", handlers.GeraGabaritoJSON)
	mux.HandleFunc("/gabaritocsv", handlers.GeraGabaritoCSV)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
