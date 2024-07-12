package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/PCPedroso/pcp-pcp-word-gen/internal/database"
	"github.com/PCPedroso/pcp-pcp-word-gen/utils"
	"github.com/tyler-smith/go-bip39"
)

var (
	value      []string
	dataSource []database.Data
	mtx        sync.Mutex
)

func GeraGabaritoJSON(w http.ResponseWriter, r *http.Request) {
	var listaWord = bip39.GetWordList()

	var gabarido []database.Gabarito
	for i, word := range listaWord {
		id := i + 1
		value := utils.InteiroParaBinario(id)
		gabarido = append(gabarido, database.Gabarito{Id: id, Word: word, Value: value})
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gabarido)
}

func GeraGabaritoCSV(w http.ResponseWriter, r *http.Request) {
	listaWord := bip39.GetWordList()
	csvWriter := csv.NewWriter(w)

	err := csvWriter.Write([]string{"Id", "Word", "Value"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, word := range listaWord {
		id := i + 1
		value := utils.InteiroParaBinario(id)

		err := csvWriter.Write([]string{strconv.Itoa(id), word, value})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	csvWriter.Flush()
	w.Header().Set("Content-Type", "text/csv")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, value)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		formData := r.FormValue("value")

		idx := formData[:1]
		vlr, err := strconv.Atoi(formData[1:])
		if err != nil {
			http.Error(w, "Dados inválidos para a operação", http.StatusMethodNotAllowed)
		}

		data := database.Data{Indice: idx, Value: vlr}

		mtx.Lock()
		value = append(value, formData)
		dataSource = append(dataSource, data)
		mtx.Unlock()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func ProcessHandler(w http.ResponseWriter, r *http.Request) {
	entropy, _ := bip39.NewEntropy(256)

	var words string
	for range 10 {
		mnemonic, _ := bip39.NewMnemonic(entropy)
		words += mnemonic + " "
	}

	lista := utils.TextToList(words)

	result := []database.Result{}
	for i, item := range dataSource {
		result = append(result, database.Result{Word: strings.Fields(utils.ReplaceWordsInt(lista[i], item.Indice, item.Value))})
	}

	tmpl := template.Must(template.ParseFiles("result.html"))
	tmpl.Execute(w, result)
}

func ClearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		mtx.Lock()
		value = []string{}
		dataSource = []database.Data{}
		mtx.Unlock()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
