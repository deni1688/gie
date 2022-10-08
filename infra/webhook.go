package infra

import (
	"deni1688/gitissue/domain"
	"fmt"
	"net/http"
)

type webhookNotifier struct {
	webhooks []string
	client   *http.Client
}

func (r webhookNotifier) Notify(issues *[]domain.Issue) error {
	for _, webhook := range r.webhooks {
		fmt.Printf("Sending to webhook: %s\n", webhook)
	}
	return nil
}

func NewWebhookNotifier(webhooks []string) domain.Notifier {
	return &webhookNotifier{webhooks, http.DefaultClient}
}
