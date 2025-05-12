package notifier_test

import (
	"errors"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service/notifier"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service/notifier/notifierfakes"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/lateralusd/textbelt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ErrStatus     = errors.New("status error")
	ErrSendFailed = errors.New("send failed")
	ErrQuotaCheck = errors.New("quota check error")
)

var _ = Describe("SmsNotifier", func() {
	var (
		fakeClient  *notifierfakes.FakeTextbeltClient
		smsNotifier *notifier.SmsNotifier
		testNotif   model.Notification
		errAction   error
	)

	BeforeEach(func() {
		fakeClient = &notifierfakes.FakeTextbeltClient{}
		smsNotifier = notifier.NewSmsNotifier(fakeClient)

		testNotif = model.Notification{
			Channel:   "sms",
			Recipient: "+359888123456",
			Subject:   "",
			Message:   "Test SMS",
		}
	})

	JustBeforeEach(func() {
		errAction = smsNotifier.Send(testNotif)
	})

	Context("when sending SMS is successful", func() {
		BeforeEach(func() {
			fakeClient.QuotaReturns(1, nil)
			fakeClient.SendReturns("msg-id", nil)
			fakeClient.StatusReturns(textbelt.StatusSent, nil)
		})

		It("succeeds", func() {
			Expect(errAction).NotTo(HaveOccurred())
		})
	})

	Context("when checking quota fails", func() {
		BeforeEach(func() {
			fakeClient.QuotaReturns(0, ErrQuotaCheck)
		})

		It("returns an error", func() {
			Expect(errAction).To(MatchError(ErrQuotaCheck))
		})
	})

	Context("when quota is zero", func() {
		BeforeEach(func() {
			fakeClient.QuotaReturns(0, nil)
		})

		It("does not send the message and returns no error", func() {
			Expect(errAction).NotTo(HaveOccurred())
			Expect(fakeClient.SendCallCount()).To(BeZero())
		})
	})

	Context("when sending the message fails", func() {
		BeforeEach(func() {
			fakeClient.QuotaReturns(1, nil)
			fakeClient.SendReturns("", ErrSendFailed)
		})

		It("returns an error", func() {
			Expect(errAction).To(MatchError(ErrSendFailed))
		})
	})

	Context("when checking message status fails", func() {
		BeforeEach(func() {
			fakeClient.QuotaReturns(1, nil)
			fakeClient.SendReturns("", nil)
			fakeClient.StatusReturns("", ErrStatus)
		})
		It("returns an error", func() {
			Expect(errAction).To(MatchError(ErrStatus))
		})
	})

	Context("when message status is failed", func() {
		BeforeEach(func() {
			fakeClient.QuotaReturns(1, nil)
			fakeClient.SendReturns("", nil)
			fakeClient.StatusReturns(textbelt.StatusFailed, nil)
		})

		It("returns an error", func() {
			Expect(errAction).To(MatchError(notifier.ErrFailedToSendSMS))
		})
	})
})
