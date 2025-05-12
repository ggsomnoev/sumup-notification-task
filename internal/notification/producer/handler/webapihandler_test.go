package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/ggsomnoev/sumup-notification-task/internal/notification/model"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/producer/handler"
	"github.com/ggsomnoev/sumup-notification-task/internal/notification/producer/handler/handlerfakes"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ErrPublishFailed = errors.New("failed to publish notification")

var _ = Describe("Notification Handler", func() {
	var (
		e         *echo.Echo
		ctx       context.Context
		recorder  *httptest.ResponseRecorder
		publisher *handlerfakes.FakePublisher
	)

	BeforeEach(func() {
		e = echo.New()
		ctx = context.Background()
		recorder = httptest.NewRecorder()
		publisher = &handlerfakes.FakePublisher{}

		handler.RegisterHandlers(ctx, e, publisher)
	})

	Describe("POST /notifications", func() {
		var (
			notification model.Notification
			req          *http.Request
		)

		BeforeEach(func() {
			notification = model.Notification{
				Channel:   "email",
				Recipient: "test@mest.com",
				Message:   "Hello",
				Subject:   "Test Subject",
			}

			body, _ := json.Marshal(notification)
			req = httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		})

		JustBeforeEach(func() {
			e.ServeHTTP(recorder, req)
		})

		It("succeeds", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(publisher.PublishCallCount()).To(Equal(1))
			_, actualMessage := publisher.PublishArgsForCall(0)
			Expect(actualMessage.UUID).NotTo(Equal(uuid.Nil))
			Expect(actualMessage.Notification).To(Equal(notification))
		})

		When("invalid JSON is posted", func() {
			BeforeEach(func() {
				req = httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewBufferString("{invalid json"))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			})

			It("returns 400", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("validation fails", func() {
			BeforeEach(func() {
				notification = model.Notification{
					Channel:   "invalid",
					Recipient: "",
					Message:   "",
				}

				body, _ := json.Marshal(notification)
				req = httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			})

			It("returns 400", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("publish fails", func() {
			BeforeEach(func() {
				publisher.PublishReturns(ErrPublishFailed)
			})

			It("returns 500", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
