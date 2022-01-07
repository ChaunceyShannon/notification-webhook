package main

import (
	. "github.com/ChaunceyShannon/golanglibs"
	"github.com/gin-gonic/gin"
)

// map = {provider name: function name}
var webHookHandlers = map[string]func(c *gin.Context) string{
	"huaweiCloud":   huaweiCloudWebhook,
	"fluxcdGeneric": fluxcdGenericWebhook,
	"httpPostRaw":   httpPostRawWebhook,
	"httpPostJSON":  httpPostJSONWebhook,
}

// map = {provider name: handler name}
var messageHandlers = map[string]messageHandler{
	"telegramBot":  &telegramMessageHandler{},
	"slackWebhook": &slackWebhookMessageHandler{},
	"slackBot":     &slackBotMessageHandler{},
	"matrixUser":   &matrixUserMessageHandler{},
	"email":        &emailMessageHandler{},
}

func main() {
	cfg := Argparser("").ParseArgs().Cfg

	Lg.SetLevel(cfg.Section("system").Key("level").Value())

	// chan for the message to send
	msgchan := make(chan messageStruct)

	ws := webhookServer{}
	ws.init(msgchan, cfg.Section("system").Key("bind").Value())

	ms := messageSender{}
	ms.init(msgchan, Int(cfg.Section("system").Key("retry").Value()))

	Lg.Trace("Read the config file")
	for _, section := range cfg.SectionStrings() {
		provider := cfg.Section(section).Key("provider").Value()

		// If it is a webhook handler, register it
		if Map(webHookHandlers).Has(provider) {
			Lg.Trace("Register webhook handler \""+provider+"\" with path:", cfg.Section(section).Key("path").Value())
			ws.register(
				cfg.Section(section).Key("path").Value(),
				cfg.Section(section).Key("messageHandlers").Value(),
				webHookHandlers[provider],
			)
		}

		// If it is a message handler, register it
		// Lg.Trace("messageHandlers keys:", Map(messageHandlers).Keys())
		// Lg.Trace("provider:", provider)
		if Map(messageHandlers).Has(provider) {
			Lg.Trace("Register message handler:", provider)
			args := make(map[string]string)
			for _, k := range cfg.Section(section).KeyStrings() {
				args[k] = cfg.Section(section).Key(k).Value()
			}
			ms.register(
				section,
				args,
				messageHandlers[provider],
			)
		}
	}

	// Run
	go ws.run()
	ms.run()
}
