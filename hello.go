package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const quantidadeMonitoramentos = 3
const delayEntreMonitoramentos = 3 * time.Second

func main() {
	exibeIntroducao()

	for {
		fmt.Println("")
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa.")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando.")
			os.Exit(-1)
		}
	}

}

func exibeIntroducao() {
	nome := "Douglas"
	versao := 1.1
	fmt.Println("Olá, sr.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func leComando() int {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair Programa")

	var comandoLido int
	fmt.Scan(&comandoLido)
	return comandoLido
}

func iniciarMonitoramento() {
	sites := leSitesDoArquivo()

	for i := 0; i < quantidadeMonitoramentos; i++ {
		fmt.Println("")
		for _, site := range sites {
			testaSite(site)
		}
		time.Sleep(delayEntreMonitoramentos)
	}
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Site:", site, "Erro: ", err)
		os.Exit(-1)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", resp.Request.URL, "OK - StatusCode:", "[", resp.StatusCode, "]")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", resp.Request.URL, "Falha - StatusCode:", "[", resp.StatusCode, "]")
		registraLog(site, false)
	}
}

func leSitesDoArquivo() []string {
	sites := []string{}
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Erro ao ler o arquivo de sites", err)
		os.Exit(-1)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Erro ao ler oa linha do arquivo", err)
			os.Exit(-1)
		}
	}

	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Erro ao escrever o log", err)
		os.Exit(-1)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " -- " + site + "- online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLogs() {
	arquivo, err := os.Open("logs.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de logs", err)
		os.Exit(-1)
	}

	conteudo, err := io.ReadAll(arquivo)

	if err != nil {
		fmt.Println("Erro ao ler o arquivo de logs", err)
		os.Exit(-1)
	}

	fmt.Println(string(conteudo))

	arquivo.Close()
}
