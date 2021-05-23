package main

import (
	"fmt"
	"sync"
	"time"
	"net/http"
	"net/url"
)

type availabilityChecker struct {
	wg sync.WaitGroup
	goroutinesLimitChan chan int
	notAvailableSites []string
	config config
}

func (ac *availabilityChecker) startChecker(urls []string) {
	for {
		for _, url := range urls {
			ac.goroutinesLimitChan <- 1
			ac.wg.Add(1)
			go ac.checkURL(url)
		}

		ac.wg.Wait()

		fmt.Println(ac.notAvailableSites)
		if len(ac.notAvailableSites) != 0 {
			ac.sendReport()
		}

		time.Sleep(ac.config.CheckerIntervalMin * time.Minute)
	}
}

func (ac *availabilityChecker) sendReport() error {
	message := "Сайты, не прошедшие проверку:"
	for _, url := range ac.notAvailableSites {
		message += url + "\n"
	}

	data := url.Values{
		"chat_id": {ac.config.TelegramChatID},
		"text": {message},
	}

	telegramMessageURL := fmt.Sprintf(ac.config.TelegramURLTemplate, ac.config.TelegramBotToken)

	_, err := http.PostForm(telegramMessageURL, data)
	if err != nil {
		return err
	}

	return nil
}

func (ac *availabilityChecker) checkURL(url string) {
	var availability bool
	client := http.Client{
		Timeout: ac.config.HttpTimeoutSec * time.Second,
	}

	resp, err := client.Get(url)
	if err == nil && resp.StatusCode == 200 {
		availability = true
	}

	if (!availability) {
		ac.notAvailableSites = append(ac.notAvailableSites, url)
	}

	ac.wg.Done()
	<- ac.goroutinesLimitChan
}