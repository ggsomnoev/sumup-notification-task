package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ggsomnoev/sumup-notification-task/internal/logger"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
)

var SlackTimeout = 5 * time.Second

type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: SlackTimeout},
	}
}

func (sn *SlackNotifier) Send(n model.Notification) error {
	payload := map[string]string{
		"text": fmt.Sprintf("*%s*\n%s", n.Subject, n.Message),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %w", err)
	}

	req, err := http.NewRequest("POST", sn.webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create Slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := sn.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack request: %w", err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.GetLogger().Errorf("could not close response body: %v", err)
		}
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("received non-2xx response from Slack: %s", response.Status)
	}

	return nil
}
