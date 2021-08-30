package main

import "github.com/gin-gonic/gin"

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
	cfg := argparser("").parseArgs().cfg

	lg.setLevel(cfg.Section("system").Key("level").Value())

	// chan for the message to send
	msgchan := make(chan messageStruct)

	ws := webhookServer{}
	ws.init(msgchan, cfg.Section("system").Key("bind").Value())

	ms := messageSender{}
	ms.init(msgchan, toInt(cfg.Section("system").Key("retry").Value()))

	lg.trace("Read the config file")
	for _, section := range cfg.SectionStrings() {
		provider := cfg.Section(section).Key("provider").Value()

		// If it is a webhook handler, register it
		if keyInMap(provider, webHookHandlers) {
			lg.trace("Register webhook handler \""+provider+"\" with path:", cfg.Section(section).Key("path").Value())
			ws.register(
				cfg.Section(section).Key("path").Value(),
				cfg.Section(section).Key("messageHandler").Value(),
				webHookHandlers[provider],
			)
		}

		// If it is a message handler, register it
		if keyInMap(provider, messageHandlers) {
			lg.trace("Register message handler:", provider)
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
