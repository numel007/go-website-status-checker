package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

func calculate(num1 int, num2 int) int {
	return (num1 * num2)
}

func printStatus(url string, statusCode int, statusText string) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	switch statusCode {
	case 200:
		println(blue("GET ") + url + " " + green(200) + " " + green(statusText))
	case 408:
		println(blue("GET ") + url + " " + red(408) + " " + red(statusText))
	default:
		println(blue("GET ") + url + " " + yellow(statusCode) + " " + yellow(statusText))
	}
}

func writeStatus(url string, statusCode int, statusText string) {
	file, err := os.OpenFile("website-status.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	logger.Println("GET " + url + " " + fmt.Sprint(statusCode) + " " + statusText)
}

func checkStatus(urls []string) {
	var wg sync.WaitGroup

	// Temporarily removed infinite loop for testing
	// for {
	wg.Add(len(urls))

	for i := 0; i < len(urls); i++ {
		go func(url string) {
			defer wg.Done()
			client := http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(url)

			if err != nil {
				writeStatus(url, 408, http.StatusText(408))
				printStatus(url, 408, http.StatusText(408))
			} else {
				if resp.StatusCode == 200 {
					writeStatus(url, 200, http.StatusText(200))
					printStatus(url, 200, http.StatusText(200))
				} else {
					writeStatus(url, resp.StatusCode, http.StatusText(resp.StatusCode))
					printStatus(url, resp.StatusCode, http.StatusText(resp.StatusCode))
				}
			}
		}(urls[i])
	}

	wg.Wait()
	// time.Sleep(1 * time.Minute)
	// }
}

func main() {
	content, err := ioutil.ReadFile("websites.txt")
	if err != nil {
		panic(err)
	}

	websites := strings.Split(string(content), "\n")
	checkStatus(websites)
}
