package notifier_test

import (
	"errors"
	"net/http"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service/notifier"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/sendgrid/rest"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service/notifier/notifierfakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ErrMailClient = errors.New("mail client failure")

var _ = Describe("EmailNotifier", func() {
	var (
		fakeClient    *notifierfakes.FakeSendGridClient
		emailNotifier *notifier.EmailNotifier
		testMessage   model.Notification
		mockResponse  *rest.Response
		errAction     error
		sender        string
	)

	BeforeEach(func() {
		sender = "from@example.com"
		fakeClient = &notifierfakes.FakeSendGridClient{}
		emailNotifier = notifier.NewEmailNotifier(fakeClient, sender)

		testMessage = model.Notification{
			Channel:   "email",
			Recipient: "to@example.com",
			Subject:   "Test Subject",
			Message:   "Hello there ;)!",
		}

		mockResponse = &rest.Response{
			StatusCode: http.StatusAccepted,
			Body:       "Accepted",
		}
	})

	JustBeforeEach(func() {
		errAction = emailNotifier.Send(testMessage)
	})

	Context("email is sent successfully", func() {
		BeforeEach(func() {
			fakeClient.SendReturns(mockResponse, nil)
		})

		It("succeeds", func() {
			Expect(errAction).NotTo(HaveOccurred())
			Expect(fakeClient.SendCallCount()).To(Equal(1))
		})
	})

	Context("there is an error", func() {
		BeforeEach(func() {
			fakeClient.SendReturns(nil, ErrMailClient)
		})

		It("returns the wrapped error", func() {
			Expect(errAction).To(MatchError(ErrMailClient))
		})
	})

	Context("the mail client returns a non-success status code", func() {
		BeforeEach(func() {
			fakeClient.SendReturns(&rest.Response{
				StatusCode: http.StatusBadRequest,
				Body:       "Bad Request",
			}, nil)
		})

		It("returns an error containing the status code", func() {
			Expect(errAction).To(MatchError(notifier.ErrFailedToSendEmail))
		})
	})
})
