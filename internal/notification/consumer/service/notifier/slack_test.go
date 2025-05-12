package notifier_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service/notifier"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SlackNotifier", func() {
	var (
		server      *httptest.Server
		notifierSvc *notifier.SlackNotifier
		requestBody []byte
	)

	BeforeEach(func() {
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]string
			data, _ := io.ReadAll(r.Body)
			requestBody = data
			_ = json.Unmarshal(data, &body)

			Expect(body["text"]).To(ContainSubstring("Subject"))
			Expect(body["text"]).To(ContainSubstring("Message"))
			w.WriteHeader(http.StatusOK)
		}))

		notifierSvc = notifier.NewSlackNotifier(server.URL)
	})

	AfterEach(func() {
		server.Close()
	})

	It("sends a Slack notification successfully", func() {
		err := notifierSvc.Send(model.Notification{
			Subject: "Subject",
			Message: "Message",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(string(requestBody)).To(ContainSubstring("Subject"))
	})

	It("returns error on non-2xx status", func() {
		server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Slack error", http.StatusBadRequest)
		})

		err := notifierSvc.Send(model.Notification{
			Subject: "Oops",
			Message: "Error",
		})
		Expect(err).To(MatchError(ContainSubstring("received non-2xx response from Slack")))
	})

	It("returns error on client failure", func() {
		notifier.SlackTimeout = 100 * time.Millisecond
		badNotifier := notifier.NewSlackNotifier("http://invalid.host")

		err := badNotifier.Send(model.Notification{
			Subject: "Fail",
			Message: "Client error",
		})
		Expect(err).To(MatchError(ContainSubstring("failed to send Slack request")))
	})
})
