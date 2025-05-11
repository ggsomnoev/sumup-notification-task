# ADR-001: SMS Service Provider Selection

## Status
Accepted

## Context
We need to select a service provider for sending SMS messages as part of the notification system. Our priorities include:

- Easy integration with Go
- To cover at least BG and/or PH
- Free tier trial for testing

## Considered Options

### 1. Twilio
- ✅ Provides both SMS and email services
- ✅ Easy to integrate (has Go SDK and solid documentation)
- ❌ Free tier only works with verified US (Australia or Ireland) phone numbers. Selecting another phone number costs a lot.
- ❌ Fails to send SMS to PH numbers under free tier

### 2. AWS SNS (Simple Notification Service)
- ✅ Free tier available
- ✅ Global support
- ❌ Requires credit card information for signup

### 3. Textlocal
- ✅ Free tier exists
- ❌ Requires a UK phone number to activate the account

### 4. Textbelt
- ✅ Free tier allows 1 SMS per day
- ✅ Works with BG (Bulgaria) numbers
- ❌ Does not support PH (Philippines) numbers due to prior abuse
- ❌ Limited quota for development purposes (1 SMS per day)

## Decision
We will initially experiment with **Textbelt** for testing with BG numbers, and **Twilio** for email and potentially US-based SMS use (testing with a US friend).

## Considerations
- Fallback logic in the service in case of provider failure.

