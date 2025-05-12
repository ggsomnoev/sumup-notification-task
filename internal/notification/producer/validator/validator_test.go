package validator_test

import (
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/producer/validator"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	var (
		n         model.Notification
		errAction error
	)

	BeforeEach(func() {
		n = model.Notification{
			Subject: "Test Message",
			Message: "Hello",
		}
	})

	JustBeforeEach(func() {
		errAction = validator.ValidateNotification(n)
	})

	When("the notification is email", func() {
		BeforeEach(func() {
			n.Channel = model.ChannelEmail
			n.Recipient = "test@mest.com"
		})

		It("succeeds", func() {
			Expect(errAction).To(Succeed())
		})

		Context("rejects invalid emails", func() {
			BeforeEach(func() {
				n.Recipient = "invalid-email"
			})

			It("returns an error", func() {
				Expect(errAction).To(MatchError(validator.ErrInvalidEmailFormat))
			})
		})

		Context("rejects missing subject field", func() {
			BeforeEach(func() {
				n.Subject = ""
			})

			It("returns an error", func() {
				Expect(errAction).To(MatchError(validator.ErrMissingFields))
			})
		})

		Context("rejects missing message field", func() {
			BeforeEach(func() {
				n.Message = ""
			})

			It("returns an error", func() {
				Expect(errAction).To(MatchError(validator.ErrMissingFields))
			})
		})
	})

	When("the notification is SMS", func() {
		BeforeEach(func() {
			n.Channel = model.ChannelSMS
			n.Recipient = "+12345678901"
		})

		It("succeeds", func() {
			Expect(errAction).To(Succeed())
		})

		Context("rejects invalid phone numbers", func() {
			BeforeEach(func() {
				n.Recipient = "not-a-phone"
			})

			It("returns an error", func() {
				Expect(errAction).To(MatchError(validator.ErrInvalidPhoneNumFormat))
			})
		})

		Context("rejects missing message field", func() {
			BeforeEach(func() {
				n.Message = ""
			})

			It("returns an error", func() {
				Expect(errAction).To(MatchError(validator.ErrMissingFields))
			})
		})
	})

	When("the notification is Slack", func() {
		BeforeEach(func() {
			n.Channel = model.ChannelSlack
		})

		It("succeeds", func() {
			Expect(errAction).To(Succeed())
		})

		Context("rejects missing subject field", func() {
			BeforeEach(func() {
				n.Subject = ""
			})

			It("returns an error", func() {
				Expect(errAction).To(MatchError(validator.ErrMissingFields))
			})
		})

		Context("rejects missing message field", func() {
			BeforeEach(func() {
				n.Message = ""
			})

			It("returns an error", func() {
				Expect(errAction).To(MatchError(validator.ErrMissingFields))
			})
		})
	})

	Context("when given unsupported channel", func() {
		var unsupportedChannel model.ChannelType = "Skype"
		
		BeforeEach(func() {
			n.Channel = unsupportedChannel
		})

		It("returns error", func() {
			Expect(errAction).To(MatchError(validator.ErrUnsupportedChannel))
		})
	})
})
