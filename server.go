package main

import (
	"net/http"
	"strconv"
	"sync"
)

type Data struct {
	Value string
}

type DataBase struct {
	Id   int
	Word string
}

var (
	dataSource []DataBase
	mtx        sync.Mutex
)

func smain() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

}

func (data *Data) AddData() (*DataBase, error) {
	id, err := strconv.Atoi(data.Value[:1])
	if err != nil {
		return nil, err
	}

	return &DataBase{
		Id:   id,
		Word: data.Value[1:],
	}, nil
}
