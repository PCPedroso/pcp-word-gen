package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

type Dados struct {
	Id    string
	Valor string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var dados []Dados

	for {
		fmt.Print("Digite dados (ou 'ok' para concluir): ")

		scanner.Scan()
		entrada := scanner.Text()

		if entrada == "ok" {
			break
		}

		if entrada == "lista" {
			nomeArquivo := "gabarito.csv"
			arquivo, err := os.Create(nomeArquivo)
			if err != nil {
				panic(err)
			}
			defer arquivo.Close()

			var listaWord = bip39.GetWordList()

			for i, word := range listaWord {
				_, err := arquivo.Write([]byte(fmt.Sprintf("%v;%v;%v;\n", i+1, word, inteiroParaBinario(i+1))))
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("Arquivo gerado com sucesso: ", nomeArquivo)
			os.Exit(0)
			break
		}

		partes := strings.Split(entrada, ";")
		dados = append(dados, Dados{Id: partes[0], Valor: partes[1]})
	}

	entropy, _ := bip39.NewEntropy(256)

	var palavras string
	for range 10 {
		mnemonic, _ := bip39.NewMnemonic(entropy)
		palavras += mnemonic + " "
	}

	lista := textoToList(palavras)

	nomeArquivo := "palavras.csv"
	arquivo, err := os.Create(nomeArquivo)
	if err != nil {
		panic(err)
	}
	defer arquivo.Close()

	for i, item := range dados {
		lista[i] = substituiPalavra(lista[i], item.Id, item.Valor)
	}

	for _, item := range lista {
		_, err := arquivo.Write([]byte(fmt.Sprint(strings.ReplaceAll(item, " ", ";"), "\n")))
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Arquivo gerado com sucesso: ", nomeArquivo)
}

func substituiPalavra(valores string, indice string, valor string) string {
	i, _ := strconv.Atoi(indice)
	palavras := strings.Fields(valores)
	decimal := binaryToDecimal(valor)
	palavras[i] = getWordByIndex(decimal)
	return strings.Join(palavras, " ")
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
