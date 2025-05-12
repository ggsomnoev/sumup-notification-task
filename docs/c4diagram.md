# Notification System App – C4 Diagram

The Notification System App is designed to accept notification publishing requests and distribute them to Email, SMS, and Slack channels. It handles validation, queuing via RabbitMQ, and processing through worker consumer/s. Below is a **Container Level View** of the system.

## Container

```plantuml
@startuml
person user

rectangle NotificationApp #line:darkblue {
    node WebAPI [
        WebAPI
        ---
        Accepts requests, performs validation, publishes to queue
    ]
    node Consumer [
        NotificationConsumer
        ---
        Reads queue, processes and delivers notifications
    ]
    database DB [
        AppDB
        ---
        Stores notification data
    ]
    queue Queue [
        RabbitMQ
    ]

    WebAPI <--> Queue : Publishes messages
    Consumer <--> Queue : Consumes messages
    Consumer <--> DB : Writes/reads delivery state
}

rectangle ExternalServices {
    node EmailSvc [
        Email Service
        ---
        External SMTP provider
    ]
    node SmsSvc [
        SMS Service (TextBelt)
        ---
        1 free SMS/day limit
    ]
    node SlackAPI [
        Slack API
    ]
}

Consumer --> EmailSvc
Consumer --> SmsSvc
Consumer --> SlackAPI
user --> WebAPI : Publishes notification

@enduml
```

## Notes:
TextBelt allows only 1 free SMS/day per IP – used in the current implementation. 
Twilio implementation exists but is not currently used and properly developed.