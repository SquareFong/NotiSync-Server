package main

import "time"

type client struct {
	effective      bool
	detail         phoneDetail
	messages       []message
	needAll        bool //true: 需要上传所有，当有allMessages的post则更改
	lastActiveTime int64
	UUID           string
	newSMS         []message //json string
}

var allClients = []client{}

func findClientByUUID(UUID string, c **client) bool {
	for i, _ := range allClients {
		if allClients[i].UUID == UUID {
			*c = &allClients[i]
			return true
		}
	}
	return false
}

func activateClient(UUID string) {
	var c *client = nil
	if findClientByUUID(UUID, &c) {
		c.lastActiveTime = time.Now().Unix()
	} else {
		var tmp = client{
			effective:      true,
			detail:         phoneDetail{},
			messages:       []message{},
			needAll:        true,
			lastActiveTime: 0,
			UUID:           UUID,
			newSMS:         []message{}}
		c = &tmp
		allClients = append(allClients, *c)
	}
}

//当get的是detail/message/allMessage/,检查是否有UUID
func getDetail(UUID string) *phoneDetail {
	activateClient(UUID)
	var c *client
	if findClientByUUID(UUID, &c) {
		return &c.detail
	}
	return nil
}

func getMessages(UUID string) *[]message {
	activateClient(UUID)
	var c *client
	if findClientByUUID(UUID, &c) {
		m := c.messages
		c.messages = []message{}
		return &m
	}
	return nil
}

func getNewMessages(UUID string) *[]message {
	activateClient(UUID)
	var c *client
	if findClientByUUID(UUID, &c) {
		m := c.newSMS
		c.newSMS = []message{}
		return &m
	}
	return nil
}

func setDetail(uuid string, detail phoneDetail) bool {
	var c *client
	if findClientByUUID(uuid, &c) {
		c.detail = detail
		return true
	}
	return false
}

func setMessages(uuid string, msg []message) bool {
	var c *client
	if findClientByUUID(uuid, &c) {
		for _,m := range msg {
			c.messages = append(c.messages, m)
		}
		return true
	}
	return false
}

func setAllMessage(uuid string, msg []message) bool {
	var c *client
	if findClientByUUID(uuid, &c) {
		if c.needAll {
			c.messages = msg
			c.needAll = false
			return true
		}
	}
	return false
}

func setNewSMS(uuid string, msg message) bool {
	var c *client
	if findClientByUUID(uuid, &c) {
		c.newSMS = append(c.newSMS, msg)
		return true
	}
	return false
}
