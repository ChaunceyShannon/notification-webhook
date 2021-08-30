package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
)

// huawei cloud webhook
func huaweiCloudWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	panicerr(err)

	j := getXPathJson(str(body))
	msg := j.first("//message").text()
	j = getXPathJson(msg)
	msg = j.first("//sms_content").text()

	return msg
}

// fluxCD generic webhook
func fluxcdGenericWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	panicerr(err)

	j := getXPathJson(str(body))
	msg := "Severity:   " + j.first("//severity").text() + "\n"
	msg += "Timestamp:  " + j.first("//timestamp").text() + "\n"
	msg += "Reason:     " + j.first("//reason").text() + "\n"
	msg += "Controller: " + j.first("//reportingController").text() + "\n"
	msg += "Object:\n"
	msg += "  Kind: " + j.first("//involvedObject/kind").text() + "\n"
	msg += "  Namespace: " + j.first("//involvedObject/namespace").text() + "\n"
	msg += "  Name: " + j.first("//involvedObject/name").text() + "\n\n"
	msg += j.first("//message").text()

	return msg
}

// http post raw webhook
func httpPostRawWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	panicerr(err)

	return str(body)
}

// http post json
func httpPostJSONWebhook(c *gin.Context) string {
	body, err := ioutil.ReadAll(c.Request.Body)
	panicerr(err)

	if json.Valid(body) {
		b, err := yaml.JSONToYAML(body)
		if err == nil {
			body = b
		}
	}

	return str(body)
}

// Another webhooks
