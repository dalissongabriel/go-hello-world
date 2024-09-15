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

const amountSiteCheck = 3
const delayBetweenChecks = 3 * time.Second

func main() {
	showWelcome()

	for {
		fmt.Println("")
		action := readAction()

		switch action {
		case 1:
			startCheck()
		case 2:
			showLogs()
		case 0:
			fmt.Println("Saindo do programa.")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando.")
			os.Exit(-1)
		}
	}

}

func showWelcome() {
	nome := "Alisson"
	fmt.Println("Olá, sr.", nome)
}

func readAction() int {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair Programa")

	var chosenOption int
	fmt.Scan(&chosenOption)
	return chosenOption
}

func startCheck() {
	sites := getSitesFromFile()

	for i := 0; i < amountSiteCheck; i++ {
		fmt.Println("")
		for _, site := range sites {
			checkSite(site)
		}
		if i != amountSiteCheck-1 {
			time.Sleep(delayBetweenChecks)
		}
	}
}

func checkSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Site:", site, "Erro: ", err)
		os.Exit(-1)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", resp.Request.URL, "OK - StatusCode:", "[", resp.StatusCode, "]")
		writeLogs(site, true)
	} else {
		fmt.Println("Site:", resp.Request.URL, "Falha - StatusCode:", "[", resp.StatusCode, "]")
		writeLogs(site, false)
	}
}

func getSitesFromFile() []string {
	sites := []string{}
	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Erro ao ler o arquivo de sites", err)
		os.Exit(-1)
	}

	reader := bufio.NewReader(file)

	for {
		row, err := reader.ReadString('\n')
		row = strings.TrimSpace(row)
		sites = append(sites, row)

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Erro ao ler oa linha do arquivo", err)
			os.Exit(-1)
		}
	}

	file.Close()

	return sites
}

func writeLogs(site string, status bool) {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Erro ao escrever o log", err)
		os.Exit(-1)
	}

	file.WriteString(time.Now().Format("02/01/2006 15:04:05") + " -- " + site + "- online: " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func showLogs() {
	file, err := os.Open("logs.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de logs", err)
		os.Exit(-1)
	}

	content, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("Erro ao ler o arquivo de logs", err)
		os.Exit(-1)
	}

	fmt.Println(string(content))

	file.Close()
}
