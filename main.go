package main

import(
	"fmt"
	"os"
	"bufio"
)

func main() {
	var ac availabilityChecker
	err := ac.config.fillFromYML("conf.yml")
	if err != nil {
		fmt.Println("Ошибка чтения файла конфигурации")
		return
	}

	file, err := os.Open(ac.config.URLsFile)
	_ = file
	if err != nil {
		fmt.Println("Ошибка чтения файла " + ac.config.URLsFile)
		return
	}

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	file.Close()

	ac.goroutinesLimitChan = make(chan int, ac.config.GoroutinesMaxCount)
	ac.startChecker(urls)
}