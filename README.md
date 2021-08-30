# notification-webhook

It is a forwarder that helps you forward the message from webhook to telegram, matrix, email, slack, etc. 

# Overview

When I am working, I am responsible for maintaining FluxCD, Flagger, Alertmanager, Azure, AWS, HuaweiClound, Jenkins, Gitea, and other services. And I am responsible for maintaining some applications that I wrote also. For example, an application that sends the message out if the log level is ERROR.

All the systems have some way to alert, but they do not all support telegram and matrix since our company is in heavy use.

So I decided to write a forwarder to do so. It will help me forward the message to telegram and matrix, and other services. For example, Slack, Email.   

# Build

```bash
$ git clone https://github.com/ChaunceyShannon/notification-webhook
$ cd notification-webhook
$ go build .
```

# Docker

```bash
$ docker run \
  -d \
  --name notification-webhook \
  -v `pwd`/kustomize/notification-webhook.ini:/app/notification-webhook.ini \
  chaunceyshannon/notification-webhook:latest 
```

# Kubernetes 

Install notification-webhook to `default` namespace

```bash
$ kubectl apply -k kustomize -n default
```

# Configuration

To use this application, you have to define at least a message handler and a webhook handler. 

First, let's define a message with the name `telegramBot-production`. It will send a message through the telegram Bot API. 

```ini
; Custom name
[telegramBot-production]
; Provider's name for outgoing message. Different provider with different arguments.
provider = telegramBot
; Telegram Bot's Token
token    = 361406910:d9Ud2O6llxNUGQxSHH0MoXHMKXoiLcYefFk
; Chat ID to send message to
chatID   = -95780003 ; production notification channel
```

And then reference it in webhook handler by name `telegramBot-production`. 

```ini
[httpPostRaw]
; The body of the post request will be sent as a message
; Provider's name for incoming webhook
provider          = httpPostRaw
; URL for the webhook, custom path
path              = /httpPostRaw-Production-URL
; Message handlers name list, separated by a comma, will send a message with these message handlers.
messageHandlers = telegramBot-production
```

Test 

```bash
$ curl -X POST https://notification-webhook.example.com/httpPostRaw-Production-URL -d "Hello World!"
```

Now you should receive the message `Hello World!` inside the telegram channel.

For complete configuration of the config file, please have a look at `kustomize/notification-webhook.ini`

# More

My Blog: https://shareitnote.com/