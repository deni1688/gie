package infra

import (
	"deni1688/gogie/internal/issues"
	"fmt"
)

type webhookNotifier struct {
	webhooks []string
	client   HttpClient
}

func (r webhookNotifier) Notify(issues *[]issues.Issue) error {
	for _, webhook := range r.webhooks {
		fmt.Printf("Sending to webhook: %s\n", webhook)
	}
	return nil
}

func NewWebhookNotifier(webhooks []string, client HttpClient) issues.Notifier {
	return &webhookNotifier{webhooks, client}
}
