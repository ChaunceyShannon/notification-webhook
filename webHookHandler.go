package main

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/ChaunceyShannon/golanglibs"

	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
)

// huawei cloud webhook
func huaweiCloudWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	Panicerr(err)

	j := String(body).JsonXPath()
	msg := j.First("//message").Text()
	j = msg.JsonXPath()
	msg = j.First("//sms_content").Text()

	return msg.S
}

// fluxCD generic webhook
func fluxcdGenericWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	Panicerr(err)

	j := String(body).JsonXPath()
	msg := "Severity:   " + j.First("//severity").Text().S + "\n"
	msg += "Timestamp:  " + j.First("//timestamp").Text().S + "\n"
	msg += "Reason:     " + j.First("//reason").Text().S + "\n"
	msg += "Controller: " + j.First("//reportingController").Text().S + "\n"
	msg += "Object:\n"
	msg += "  Kind: " + j.First("//involvedObject/kind").Text().S + "\n"
	msg += "  Namespace: " + j.First("//involvedObject/namespace").Text().S + "\n"
	msg += "  Name: " + j.First("//involvedObject/name").Text().S + "\n\n"
	msg += j.First("//message").Text().S

	return msg
}

// http post raw webhook
func httpPostRawWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	Panicerr(err)

	Lg.Trace("Body bytes:", body)
	Lg.Trace("Body:", Str(body))
	return Str(body)
}

// http post json
func httpPostJSONWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	Panicerr(err)

	if json.Valid(body) {
		b, err := yaml.JSONToYAML(body)
		if err == nil {
			body = b
		}
	}

	return Str(body)
}

// Another webhooks
