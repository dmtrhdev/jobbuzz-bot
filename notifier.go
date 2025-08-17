package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Telegram has 4096-character limit for messages, using 4000 to be safe.
const telegramMsgLimit = 4000

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
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d new jobs found\n\n", len(jobs)))
	for _, job := range jobs {
		line := fmt.Sprintf("%s в %s (%s)\n\n", job.Title, job.Company, job.URL)
		if buf.Len()+len(line) > telegramMsgLimit {
			if err := t.send(buf.String()); err != nil {
				return err
			}
			buf.Reset()
		}
		buf.WriteString(line)
	}
	if buf.Len() > 0 {
		return t.send(buf.String())
	}
	return nil
}

func (t *TelegramBot) send(text string) error {
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

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: status %d", resp.StatusCode)
	}

	return nil
}
