package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"web_penetration/internal/model"
)

// Webhook通知器
type WebhookNotifier struct {
	URL     string
	Headers map[string]string
}

func (n *WebhookNotifier) Send(alert *model.DirScanAlertLog) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", n.URL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	for k, v := range n.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Email通知器
type EmailNotifier struct {
	SMTPConfig map[string]string
}

func (n *EmailNotifier) Send(alert *model.DirScanAlertLog) error {
	// TODO: 实现邮件发送
	return nil
}
