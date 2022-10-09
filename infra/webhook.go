package infra

import (
	"deni1688/gitissue/domain"
	"fmt"
)

type webhookNotifier struct {
	webhooks []string
	client   HttpClient
}

func (r webhookNotifier) Notify(issues *[]domain.Issue) error {
	for _, webhook := range r.webhooks {
		fmt.Printf("Sending to webhook: %s\n", webhook)
	}
	return nil
}

func NewWebhookNotifier(webhooks []string, client HttpClient) domain.Notifier {
	return &webhookNotifier{webhooks, client}
}
