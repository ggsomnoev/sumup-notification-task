package store_test

import (
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/consumer/store"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	When("created", func() {
		It("exists", func() {
			Expect(store.NewStore(nil)).NotTo(BeNil())
		})
	})

	Describe("instance", Serial, func() {
		// TODO: test each method.
	})
})
