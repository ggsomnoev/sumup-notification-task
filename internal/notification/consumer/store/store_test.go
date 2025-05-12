package store_test

import (
	"context"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/store"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Store", func() {
	When("created", func() {
		It("exists", func() {
			Expect(store.NewStore(nil)).NotTo(BeNil())
		})
	})

	var _ = Describe("instance", Serial, func() {
		var (
			s           *store.Store
			msg         model.Message
			errAction   error
			messageUUID uuid.UUID
		)

		BeforeEach(func() {
			s = store.NewStore(pool)

			messageUUID = uuid.New()
			msg = model.Message{
				UUID: messageUUID,
				Notification: model.Notification{
					Channel: model.ChannelSlack,
					Subject: "test",
					Message: "test message",
				},
			}
		})

		Describe("AddMessage", func() {
			JustBeforeEach(func() {
				errAction = s.RunInAtomically(ctx, func(ctx context.Context) error {
					return s.AddMessage(ctx, msg)
				})
			})

			JustAfterEach(func() {
				err := s.DeleteMessageByUUID(ctx, messageUUID)
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds", func() {
				Expect(errAction).NotTo(HaveOccurred())
			})

			It("inserts the correct message", func() {
				var storedMsg, err = s.GetMessageByUUID(ctx, messageUUID)
				Expect(err).NotTo(HaveOccurred())
				Expect(storedMsg).To(MatchFields(IgnoreExtras, Fields{
					"UUID": Equal(msg.UUID),
					"Notification": MatchFields(IgnoreExtras, Fields{
						"Channel":   Equal(msg.Channel),
						"Recipient": Equal(msg.Recipient),
						"Subject":   Equal(msg.Subject),
						"Message":   Equal(msg.Message),
					}),
				}))
			})
		})

		Describe("MarkCompleted", func() {
			BeforeEach(func() {
				err := s.RunInAtomically(ctx, func(ctx context.Context) error {
					return s.AddMessage(ctx, msg)
				})
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err := s.DeleteMessageByUUID(ctx, messageUUID)
				Expect(err).NotTo(HaveOccurred())
			})

			JustBeforeEach(func() {
				errAction = s.RunInAtomically(ctx, func(txCtx context.Context) error {
					return s.MarkCompleted(txCtx, msg.UUID)
				})
			})

			It("succeeds", func() {
				Expect(errAction).NotTo(HaveOccurred())
			})

			It("marks the message as completed", func() {
				completedAt, err := s.GetCompletedAtByUUID(ctx, msg.UUID)
				Expect(err).NotTo(HaveOccurred())
				Expect(completedAt).NotTo(BeZero())
			})
		})

		Describe("MessageExists", func() {
			var exists bool
			JustBeforeEach(func() {
				errAction = s.RunInAtomically(ctx, func(txCtx context.Context) error {
					var err error
					exists, err = s.MessageExists(txCtx, uuid.New())
					return err
				})
			})

			It("returns false for non-existing UUID", func() {
				Expect(exists).To(BeFalse())
			})

			Context("and a message is added", func() {
				var exists bool
				BeforeEach(func() {
					err := s.RunInAtomically(ctx, func(txCtx context.Context) error {
						return s.AddMessage(txCtx, msg)
					})
					Expect(err).NotTo(HaveOccurred())
				})

				AfterEach(func() {
					err := s.DeleteMessageByUUID(ctx, messageUUID)
					Expect(err).NotTo(HaveOccurred())
				})

				JustBeforeEach(func() {
					errAction = s.RunInAtomically(ctx, func(txCtx context.Context) error {
						var err error
						exists, err = s.MessageExists(txCtx, msg.UUID)
						return err
					})
				})

				It("succeeds", func() {
					Expect(errAction).NotTo(HaveOccurred())
				})

				It("returns true for inserted message", func() {
					Expect(exists).To(BeTrue())
				})
			})
		})
	})
})
