package webhook

import (
	"bytes"
	"deni1688/gie/adapters/shared"
	"deni1688/gie/core"
	"encoding/json"
	"fmt"
	"net/http"
)

type webhookNotifier struct {
	webhooks []string
	client   shared.HttpClient
}

func New(webhooks []string, client shared.HttpClient) core.Notifier {
	return &webhookNotifier{webhooks, client}
}

func (r webhookNotifier) Notify(issues *[]core.Issue) error {
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

	buffer := bytes.NewBuffer(payload)
	for _, webhook := range r.webhooks {
		req, err = http.NewRequest("POST", webhook, buffer)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err = r.client.Do(req)
		if err != nil || resp.StatusCode >= 400 {
			fmt.Printf("Error sending webhook %s %s\n", webhook, err)
		}
		defer resp.Body.Close()
	}

	return nil
}
