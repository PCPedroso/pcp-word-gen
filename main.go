package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/tyler-smith/go-bip39"
)

const nroItens int = 12

type Result struct {
	Word []string
}

type App struct {
	value  []string
	result []Result
}

type Dados struct {
	Id    string
	Valor string
}

var (
	values []string
	mu     sync.Mutex
)

type Gabarito struct {
	Id    int
	Word  string
	Value string
}

func main() {
	var app = &App{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/back", app.backToHome)
	mux.HandleFunc("/add", app.addValueHandler)
	mux.HandleFunc("/process", app.processValuesHandler)

	mux.HandleFunc("/listajson", geraGabaritoJSON)
	mux.HandleFunc("/listacsv", geraGabaritoCSV)
	mux.HandleFunc("/getvalor", handleGetValues)

	http.ListenAndServe(":8080", mux)
}

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, app.value)
}

func (app *App) backToHome(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, app.value)
}

func (app *App) addValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		valueStr := r.FormValue("value")
		_, err := strconv.Atoi(valueStr)
		if err == nil {
			app.value = append(app.value, valueStr)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) processValuesHandler(w http.ResponseWriter, r *http.Request) {
	entropy, _ := bip39.NewEntropy(256)

	var palavras string
	for range 10 {
		mnemonic, _ := bip39.NewMnemonic(entropy)
		palavras += mnemonic + " "
	}

	lista := textoToList(palavras)

	for i, item := range app.value {
		id, err := strconv.Atoi(item[1:])
		if err != nil {
			http.Error(w, "Valor inválido", http.StatusBadRequest)
		}

		app.result = append(app.result, Result{Word: strings.Fields(replaceWords(lista[i], item[1:], id))})
	}

	tmpl := template.Must(template.ParseFiles("result.html"))
	tmpl.Execute(w, app.result)
}

func handleGetValues(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var dados []Dados

	if len(values) < nroItens {
		http.Error(w, fmt.Sprintf("%v valores informados, são necessários %v valores", len(values), nroItens), http.StatusBadRequest)
		return
	}

	for _, item := range values {
		id, err := strconv.Atoi(item[1:])
		if err != nil {
			http.Error(w, "Valor inválido", http.StatusBadRequest)
		}

		dados = append(dados, Dados{Id: item[:1], Valor: inteiroParaBinario(id)})
	}

	values = nil

	entropy, _ := bip39.NewEntropy(256)

	var palavras string
	for range 10 {
		mnemonic, _ := bip39.NewMnemonic(entropy)
		palavras += mnemonic + " "
	}

	csvWriter := csv.NewWriter(w)
	lista := textoToList(palavras)

	for i, item := range dados {
		err := csvWriter.Write(strings.Fields(substituiPalavra(lista[i], item.Id, item.Valor)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	csvWriter.Flush()
	w.Header().Set("Context-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
}

func geraGabaritoJSON(w http.ResponseWriter, r *http.Request) {
	var listaWord = bip39.GetWordList()

	var gabarido []Gabarito
	for i, word := range listaWord {
		id := i + 1
		value := inteiroParaBinario(id)
		gabarido = append(gabarido, Gabarito{Id: id, Word: word, Value: value})
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gabarido)
}

func geraGabaritoCSV(w http.ResponseWriter, r *http.Request) {
	listaWord := bip39.GetWordList()
	csvWriter := csv.NewWriter(w)

	err := csvWriter.Write([]string{"Id", "Word", "Value"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, word := range listaWord {
		id := i + 1
		value := inteiroParaBinario(id)

		err := csvWriter.Write([]string{strconv.Itoa(id), word, value})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	csvWriter.Flush()
	w.Header().Set("Content-Type", "text/csv")
}

func substituiPalavra(valores string, indice string, valor string) string {
	i, _ := strconv.Atoi(indice)
	palavras := strings.Fields(valores)
	decimal := binaryToDecimal(valor)
	palavras[i] = getWordByIndex(decimal)
	return strings.Join(palavras, " ")
}

func replaceWords(values string, index string, decimal int) string {
	i, _ := strconv.Atoi(index)
	words := strings.Fields(values)
	words[i] = getWordByIndex(decimal)
	return strings.Join(words, " ")
}

func getWordByIndex(i int) string {
	return bip39.GetWordList()[i-1]
}

func binaryToDecimal(binario string) int {
	binario = strings.Replace(binario, " ", "", -1)
	numero, err := strconv.ParseInt(binario, 2, 64)
	if err != nil {
		panic(err)
	}

	return int(numero)
}

func inteiroParaBinario(numero int) string {
	binario := strconv.FormatInt(int64(numero), 2)

	for len(binario) < 12 {
		binario = "0" + binario
	}

	binarioAgrupado := ""
	for i := 0; i < len(binario); i += 4 {
		binarioAgrupado += binario[i:i+4] + " "
	}

	binarioAgrupado = strings.TrimRight(binarioAgrupado, " ")

	return binarioAgrupado
}

func textoToList(texto string) []string {
	palavras := strings.Fields(texto)

	registros := make([]string, 0)
	registroAtual := ""

	for _, palavra := range palavras {
		if len(strings.Fields(registroAtual)) >= 10 {
			registros = append(registros, registroAtual)
			registroAtual = ""
		}

		registroAtual += palavra + " "
	}

	if registroAtual != "" {
		registros = append(registros, registroAtual)
	}

	return registros
}
