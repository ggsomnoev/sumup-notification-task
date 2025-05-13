# Notification Sending System

## Description
This app is a tiny demonstration of sending messages to clients over different media:
* Accepts **HTTP** requests to publish notifications (Email, Slack, SMS).  
* Pushes every request onto a **RabbitMQ** queue.
* Retries messages on failures couple of times. 
* **Consumer/s** processes the queue and delivers the message through:
  * **SendGrid** – Email  
  * **Slack WebHooks** – Slack  
  * **TextBelt** – SMS  
* Persists every processed/ing message in PostgreSQL:
  * used for tracking.  

## Requirements
- Support sending notifications through:
  - SMS
  - Slack
  - Email
- Delivery - "at least once"
- HTTP API to accept messages
- Horizontaly scalable
- Easily extendable support for new channels
- Optional: queue based processing

## Architecture

The project follows [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/The-Clean-Architecture.html) to get:

* **easy unit-testing & mocking**  
* **clear separation of concerns**  
* **plug-and-play transports or providers** 

Key libraries & tools:

| Concern | Library |
|---------|---------|
| HTTP server | [Echo](https://echo.labstack.com/) |
| DB migrations | [golang-migrate/migrate](https://github.com/golang-migrate/migrate) |
| Testing | [Ginkgo](https://onsi.github.io/ginkgo/) + [Gomega](https://onsi.github.io/gomega/) |
| Fakes | [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) |
| Queue | RabbitMQ (`amqp091-go`) |

Useful links:
* [C4 Diagram](./docs/c4diagram.md)

## How to run the server

To run the server using docker:
```bash
make run-docker
```

If you want to manually run the migrations:

```bash
make migrate-down
make migrate-up
```

## How to run the tests

To run the unit tests:

```bash
make test
```

## Example requests and responses

### To publish an email notification

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
        "channel": "email",
        "recipient": "user@example.com",
        "subject": "Hello",
        "message": "Greetings from GO notification app!"
      }'
```

### To publish a SMS notification

⚠️ **Note**: The current free TextBelt service supports **only 1 SMS per day** from the same IP. For higher volume, a paid API key is required. See [https://textbelt.com](https://textbelt.com) for more info.

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
        "channel": ""sms",
        "recipient": "+359881234567",
        "message": "I hope this message finds you well!"
      }'
```

### To publish a slack notification

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Content-Type: application/json" \
  -d '{
        "channel": "slack",
        "subject": "Hello Reader",
        "message": "Greetings from GO notification app!"
      }'
```

## Environment Variables Setup Guide

### SendGrid Configuration

To enable email delivery via SendGrid, set the following environment variables in your `.env` file:

- `SENDGRID_API_KEY`
- `SENDGRID_SENDER_IDENTITY_EMAIL`

#### Steps to generate:

1. Create a SendGrid Account  
   Sign up at [https://sendgrid.com](https://sendgrid.com).

2. Generate an API Key
   - Navigate to **Settings > API Keys**
   - Click **Create API Key**
   - Provide a name and choose appropriate access (e.g., "Full Access", etc.)
   - Copy the generated API key and assign it to `SENDGRID_API_KEY`

3. Verify a Sender Identity
   - Go to **Email API > Sender Authentication**
   - Select **Single Sender Verification**
   - Add your email address and complete the verification process
   - Use the verified email as `SENDGRID_SENDER_IDENTITY_EMAIL`


### Slack Webhook Configuration

To send notifications to Slack, set the following environment variable:

- `SLACK_WEBHOOK_URL`

#### Steps to generate:

1. Create a Slack App
   - Visit [https://api.slack.com/apps](https://api.slack.com/apps)
   - Click **Create New App**

2. Enable Incoming Webhooks
   - Inside your app, go to **Incoming Webhooks**
   - Toggle to **Enable Incoming Webhooks**

3. Create a Webhook
   - Click **Add New Webhook to Workspace**
   - Select the desired channel to post messages to
   - Copy the generated Webhook URL and set it as `SLACK_WEBHOOK_URL`


## Future work:
* CI/CD & Terraform for cloud deployment.
* Prometheus metrics (/metrics) & OpenTelemetry tracing.
* Health-check endpoints.
* Retry / back-off strategy per messaging provider.
* Rate-limit outbound provider calls.
* Proper storage of secrets (Vault, k8s secrets, etc.)
* more tests...

Enjoy exploring the notification app, feedback and questions are welcome!

