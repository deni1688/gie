package infra

import (
	"bytes"
	"deni1688/gogie/internal/issues"
	"encoding/json"
	"fmt"
	"net/http"
)

type webhookNotifier struct {
	webhooks []string
	client   HttpClient
}

func NewWebhookNotifier(webhooks []string, client HttpClient) issues.Notifier {
	return &webhookNotifier{webhooks, client}
}

func (r webhookNotifier) Notify(issues *[]issues.Issue) error {
	if len(r.webhooks) < 1 || r.webhooks == nil {
		return nil
	}

	var (
		err  error
		req  *http.Request
		resp *http.Response
	)

	var payload []byte
	payload, err = json.Marshal(issues)
	if err != nil {
		return err
	}

	for _, webhook := range r.webhooks {
		req, err = http.NewRequest("POST", webhook, bytes.NewBuffer(payload))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err = r.client.Do(req)
		if err != nil || resp.StatusCode >= 400 {
			fmt.Printf("Error sending webhook %s %s\n", webhook, err)
		}
	}

	return nil
}
