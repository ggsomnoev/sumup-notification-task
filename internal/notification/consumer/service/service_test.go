package service_test

import (
	"context"

	"github.com/google/uuid"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/service/servicefakes"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		fakeStore    *servicefakes.FakeStore
		fakeNotifier *servicefakes.FakeNotifier
		svc          *service.Service
		msg          model.Message
		ctx          context.Context
		errAction    error
	)

	BeforeEach(func() {
		ctx = context.Background()
		fakeStore = &servicefakes.FakeStore{}
		fakeNotifier = &servicefakes.FakeNotifier{}

		svc = service.NewService(fakeStore, map[string]service.Notifier{
			"email": fakeNotifier,
		})

		msg = model.Message{
			UUID: uuid.New(),
			Notification: model.Notification{
				From:      "sender@example.com",
				Channel:   "email",
				Recipient: "test@example.com",
				Subject:   "Test",
				Message:   "Hello!",
			},
		}

		fakeStore.RunInAtomicallyStub = func(_ context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}
	})

	JustBeforeEach(func() {
		errAction = svc.Send(ctx, msg)
	})

	Context("the message is send successfully", func() {
		BeforeEach(func() {
			fakeStore.MessageExistsReturns(false, nil)
			fakeNotifier.SendReturns(nil)
			fakeStore.AddMessageReturns(nil)
			fakeStore.MarkCompletedReturns(nil)
		})

		It("succeeds", func() {
			Expect(errAction).NotTo(HaveOccurred())

			Expect(fakeStore.MessageExistsCallCount()).To(Equal(1))
			Expect(fakeNotifier.SendCallCount()).To(Equal(1))
			Expect(fakeStore.AddMessageCallCount()).To(Equal(1))
			Expect(fakeStore.MarkCompletedCallCount()).To(Equal(1))
		})
	})

	Context("skips sending the message if it already exists", func() {
		BeforeEach(func() {
			fakeStore.MessageExistsReturns(true, nil)
		})

		It("succeeds", func() {
			Expect(errAction).NotTo(HaveOccurred())
			Expect(fakeNotifier.SendCallCount()).To(Equal(0))
		})
	})

	Context("and no notifier is passed", func() {
		BeforeEach(func() {
			svc = service.NewService(fakeStore, map[string]service.Notifier{}) // no notifiers
		})

		It("returns an error", func() {
			Expect(errAction).To(HaveOccurred())
			Expect(errAction.Error()).To(ContainSubstring("no notifier registered"))
		})
	})
})
