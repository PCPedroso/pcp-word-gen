package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

type Results struct {
	RandomData []int
}

type Apps struct {
	values  []int
	results []Results
}

func server() {
	app := &Apps{}

	// http.HandleFunc("/", app.homeHandler)
	// http.HandleFunc("/add", app.addValueHandler)
	http.HandleFunc("/process", app.processValuesHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// func (app *Apps) homeHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl := template.Must(template.ParseFiles("index.html"))
// 	tmpl.Execute(w, app.values)
// }

// func (app *Apps) addValueHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		valueStr := r.FormValue("value")
// 		value, err := strconv.Atoi(valueStr)
// 		if err == nil {
// 			app.values = append(app.values, value)
// 		}
// 	}
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }

func (app *Apps) processValuesHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	app.results = nil // Clear previous random data

	for _, value := range app.values {
		randomValues := generateRandomValues(value)
		app.results = append(app.results, Results{RandomData: randomValues})
	}

	tmpl := template.Must(template.ParseFiles("result.html"))
	tmpl.Execute(w, app.results)
}

func generateRandomValues(n int) []int {
	randomValues := make([]int, 5)
	for i := 0; i < 5; i++ {
		randomValues[i] = rand.Intn(100) // Random values between 0 and 99
	}
	return randomValues
}
