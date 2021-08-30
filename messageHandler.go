package main

import (
	"net/smtp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type messageHandler interface {
	init(args map[string]string)
	send(msg string)
}

// Telegram message handler
type telegramMessageHandler struct {
	telegram *tgbotapi.BotAPI
	chatID   int64
}

func (c *telegramMessageHandler) init(args map[string]string) {
	var err error
	c.telegram, err = tgbotapi.NewBotAPI(args["token"])
	panicerr(err)

	c.chatID = toInt64(args["chatID"])
}

func (c *telegramMessageHandler) send(msg string) {
	_, err := c.telegram.Send(tgbotapi.NewMessage(c.chatID, msg))
	panicerr(err)
}

// Slack webhook
type slackWebhookMessageHandler struct {
	url string
}

func (c *slackWebhookMessageHandler) init(args map[string]string) {
	c.url = args["url"]
}

func (c *slackWebhookMessageHandler) send(msg string) {
	resp := httpPostJSON(c.url, map[string]string{"text": msg}, httpHeader{"Content-type": "application/json"})
	if resp.statusCode != 200 {
		lg.error("Error while sending message with slack webhook:", resp.content)
	}
}

// Slack Bot
type slackBotMessageHandler struct {
	chatID string
	token  string
}

func (c *slackBotMessageHandler) init(args map[string]string) {
	c.chatID = args["chatID"]
	c.token = args["token"]
}

func (c *slackBotMessageHandler) send(msg string) {
	resp := httpPostJSON("https://slack.com/api/chat.postMessage", map[string]string{"channel": c.chatID, "text": msg}, httpHeader{"Authorization": "Bearer " + c.token})
	if resp.statusCode != 200 {
		lg.error("Error while sending message with slack bot:", resp.content)
	}
}

// matrix user
type matrixUserMessageHandler struct {
	cli *matrixStruct
}

func (c *matrixUserMessageHandler) init(args map[string]string) {
	c.cli = getMatrix(args["server"]).setRoomID(args["roomID"])
	if keyInMap("userid", args) && keyInMap("token", args) {
		c.cli.setToken(args["userid"], args["token"])
	} else if keyInMap("username", args) && keyInMap("password", args) {
		c.cli.login(args["username"], args["password"])
	} else {
		lg.error("Error while initializing handler \"martixUser\", need (username and password) or (userid and token).")
		exit(0)
	}
}

func (c *matrixUserMessageHandler) send(msg string) {
	c.cli.send(msg)
}

// email
type emailMessageHandler struct {
	server   string
	from     string
	password string
	to       []string
}

func (c *emailMessageHandler) init(args map[string]string) {
	c.server = args["server"]
	c.from = args["from"]
	c.password = args["password"]

	var toarr []string
	for _, t := range strSplit(args["to"]) {
		toarr = append(toarr, strStrip(t))
	}
	c.to = toarr
}

func (c *emailMessageHandler) send(msg string) {
	var subject string
	if strIn("\n", msg) {
		subject = strSplitlines(msg)[0]
		msg = strJoin("\n", strSplitlines(msg)[1:])
	} else {
		subject = msg
	}

	msg = "From: " + c.from + "\n" +
		"To: " + strJoin(",", c.to) + "\n" +
		"Subject: " + subject + "\n\n" +
		msg

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", c.from, c.password, "smtp.gmail.com"),
		c.from, c.to, []byte(msg))

	panicerr(err)
}

// Another handlers
