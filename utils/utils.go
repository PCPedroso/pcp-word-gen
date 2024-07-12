package utils

import (
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip39"
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

func ReplaceWordsInt(values string, index string, decimal int) string {
	i, _ := strconv.Atoi(index)
	words := strings.Fields(values)
	words[i] = GetWordByIndex(decimal)
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
