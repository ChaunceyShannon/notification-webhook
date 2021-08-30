package main

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
		if err := try(func() {
			mh.init(args)
		}).Error; err != nil {
			lg.error("Error while initializing message handler \""+name+"\":", err, ". Retrying...")
			sleep(3)
		} else {
			break
		}
	}
	c.handlers[name] = mh
}

func (c *messageSender) run() {
	if len(c.handlers) == 0 {
		lg.error("At least one handler should be enabled.")
		exit(0)
	}

	// Call send function in handler for message
	for msg := range c.msgchan {
		for _, mhn := range msg.messageHandler {
			if keyInMap(strStrip(mhn), c.handlers) {
				for i := 0; i < c.retry; i++ {
					if err := try(func() {
						lg.trace("Sending message with handler \""+mhn+"\":", msg.message)
						c.handlers[mhn].send(msg.message)
					}).Error; err != nil {
						lg.error("Error while sending message:", err, ". Retrying..."+str(i)+"...")
						sleep(3)
					} else {
						break
					}
				}
			}
		}
	}
}
