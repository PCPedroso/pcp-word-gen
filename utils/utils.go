package utils

import (
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/PCPedroso/pcp-pcp-word-gen/internal/database"
	"github.com/tyler-smith/go-bip39"
	"gopkg.in/yaml.v3"
)

func InteiroParaBinario(numero int) string {
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

func ReplaceWordsString(valores string, indice string, valor string) string {
	i, _ := strconv.Atoi(indice)
	palavras := strings.Fields(valores)
	decimal := BinaryToDecimal(valor)
	palavras[i] = GetWordByIndex(decimal)
	return strings.Join(palavras, " ")
}

func BinaryToDecimal(binario string) int {
	binario = strings.Replace(binario, " ", "", -1)
	numero, err := strconv.ParseInt(binario, 2, 64)
	if err != nil {
		panic(err)
	}

	return int(numero)
}

func ReplaceWordsInt(values string, index int, decimal int) string {
	words := strings.Fields(values)
	if index == -1 {
		words[0] = "INDEX-ERROR"
	} else {
		words[index] = GetWordByIndex(decimal)
	}

	return strings.Join(words, " ")
}

func GetWordByIndex(i int) string {
	return bip39.GetWordList()[i-1]
}

func TextToList(texto string) []string {
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

func ReadConfig() (*database.Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return &database.Config{}, err
	}

	var config *database.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return &database.Config{}, err
	}

	return config, nil
}

func SetIndexByChar(value string) (int, error) {
	columns, err := ReadConfig()
	if err != nil {
		return -1, err
	}

	for key, item := range columns.Columns {
		chars := strings.Split(item, ",")

		if slices.Index(chars, value) > -1 {
			i, err := strconv.Atoi(key)
			if err != nil {
				return i, err
			}

			return i, nil
		}
	}
	return -1, nil
}

func CreateAllWords(userWords string, maxWord int) string {
	var allWords string

	for len(strings.Fields(allWords)) < maxWord {
		entropy, _ := bip39.NewEntropy(128)
		mnemonic, _ := bip39.NewMnemonic(entropy)

		for _, item := range strings.Split(mnemonic, " ") {
			if !strings.Contains(allWords, item) && !strings.Contains(userWords, item) {
				allWords += item + " "
			}
		}
	}

	return allWords
}
