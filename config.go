package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type config struct {
	HttpTimeoutSec time.Duration     `yaml:"HttpTimeoutSec"`
	GoroutinesMaxCount int           `yaml:"GoroutinesMaxCount"`
	URLsFile string                  `yaml:"URLsFile"`
	TelegramBotToken string          `yaml:"TelegramBotToken"`
	TelegramChatID string            `yaml:"TelegramChatID"`
	TelegramURLTemplate string       `yaml:"TelegramURLTemplate"`
	CheckerIntervalMin time.Duration `yaml:"CheckerIntervalMin"`
}

func (config *config) fillFromYML(ymlFileName string) error {
	content, err := ioutil.ReadFile(ymlFileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}