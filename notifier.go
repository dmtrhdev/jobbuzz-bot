package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type TelegramBot struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegramBot(token, chatID string) *TelegramBot {
	return &TelegramBot{
		token:  token,
		chatID: chatID,
		client: &http.Client{},
	}
}

func (t *TelegramBot) Send(jobs []Job) error {
	msg := t.formatMessage(jobs)
	return t.sendMessage(msg)
}

func (t *TelegramBot) formatMessage(jobs []Job) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d new jobs found\n\n", len(jobs)))

	for _, j := range jobs {
		b.WriteString(fmt.Sprintf("%s в %s (%s)\n\n", j.Title, j.Company, j.URL))
	}

	return b.String()
}

func (t *TelegramBot) sendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	payload := map[string]interface{}{
		"chat_id":                  t.chatID,
		"text":                     text,
		"disable_web_page_preview": true,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("telegram: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api %d", resp.StatusCode)
	}

	return nil
}
