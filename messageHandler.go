package main

import (
	"net/smtp"

	. "github.com/ChaunceyShannon/golanglibs"

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
	Panicerr(err)

	c.chatID = Int64(args["chatID"])
}

func (c *telegramMessageHandler) send(msg string) {
	_, err := c.telegram.Send(tgbotapi.NewMessage(c.chatID, msg))
	Panicerr(err)
}

// Slack webhook
type slackWebhookMessageHandler struct {
	url string
}

func (c *slackWebhookMessageHandler) init(args map[string]string) {
	c.url = args["url"]
}

func (c *slackWebhookMessageHandler) send(msg string) {
	resp := Http.PostJSON(c.url, map[string]string{"text": msg}, HttpHeader{"Content-type": "application/json"})
	if resp.StatusCode != 200 {
		Lg.Error("Error while sending message with slack webhook:", resp.Content)
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
	resp := Http.PostJSON("https://slack.com/api/chat.postMessage", map[string]string{"channel": c.chatID, "text": msg}, HttpHeader{"Authorization": "Bearer " + c.token})
	if resp.StatusCode != 200 {
		Lg.Error("Error while sending message with slack bot:", resp.Content)
	}
}

// matrix user
type matrixUserMessageHandler struct {
	cli *MatrixStruct
}

func (c *matrixUserMessageHandler) init(args map[string]string) {
	c.cli = Tools.Matrix(args["server"]).SetRoomID(args["roomID"])
	if Map(args).Has("userid") && Map(args).Has("token") {
		c.cli.SetToken(args["userid"], args["token"])
	} else if Map(args).Has("username") && Map(args).Has("password") {
		c.cli.Login(args["username"], args["password"])
	} else {
		Lg.Error("Error while initializing handler \"martixUser\", need (username and password) or (userid and token).")
		Os.Exit(0)
	}
}

func (c *matrixUserMessageHandler) send(msg string) {
	c.cli.Send(msg)
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
	for _, t := range String(args["to"]).Split() {
		toarr = append(toarr, t.Strip().S)
	}
	c.to = toarr
}

func (c *emailMessageHandler) send(msg string) {
	var subject string
	if String("\n").In(msg) {
		subject = String(msg).Splitlines()[0].S
		msg = String("\n").Join(String(msg).Splitlines()[1:]).S
	} else {
		subject = msg
	}

	msg = "From: " + c.from + "\n" +
		"To: " + String(",").Join(c.to).S + "\n" +
		"Subject: " + subject + "\n\n" +
		msg

	err := smtp.SendMail(c.server,
		smtp.PlainAuth("", c.from, c.password, c.server),
		c.from, c.to, []byte(msg))

	Panicerr(err)
}

// Another handlers
