package main

import . "github.com/ChaunceyShannon/golanglibs"

type messageSender struct {
	msgchan  chan messageStruct
	handlers map[string]messageHandler
	retry    int
}

func (c *messageSender) init(msgchan chan messageStruct, retry int) {
	c.msgchan = msgchan
	c.handlers = make(map[string]messageHandler)
	c.retry = retry
}

func (c *messageSender) register(name string, args map[string]string, mh messageHandler) {
	for {
		if err := Try(func() {
			mh.init(args)
		}).Error; err != nil {
			Lg.Error("Error while initializing message handler \""+name+"\":", err, ". Retrying...")
			Time.Sleep(3)
		} else {
			break
		}
	}
	c.handlers[name] = mh
}

func (c *messageSender) run() {
	if len(c.handlers) == 0 {
		Lg.Error("At least one handler should be enabled.")
		Os.Exit(0)
	}

	// Call send function in handler for message
	for msg := range c.msgchan {
		for _, mhn := range msg.messageHandler {
			Lg.Trace("mhn:", mhn)
			Lg.Trace("handlers:", Map(c.handlers).Keys())
			if Map(c.handlers).Has(String(mhn).Strip().S) {
				for i := 0; i < c.retry; i++ {
					if err := Try(func() {
						Lg.Trace("Sending message with handler \""+mhn+"\":", msg.message)
						c.handlers[mhn].send(msg.message)
					}).Error; err != nil {
						Lg.Error("Error while sending message:", err, ". Retrying..."+Str(i)+"...")
						Time.Sleep(3)
					} else {
						break
					}
				}
			}
		}
	}
}
