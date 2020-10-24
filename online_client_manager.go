package main

type client struct {
	effective      bool
	detail         phoneDetail
	messages       []message
	needAll        bool //true: 需要上传所有，当有allMessages的post则更改
	lastActiveTime int
	UUID           string
}

var allClients = []client{}

func findClientByUUID(UUID string, c **client) bool {
	for _, value := range allClients {
		if value.UUID == UUID {
			*c = &value
			return true
		}
	}
	return false
}

func activateClient(UUID string) {
	var c *client
	if !findClientByUUID(UUID, &c) {
		*c = client{
			effective:      true,
			detail:         phoneDetail{},
			messages:       []message{},
			needAll:        true,
			lastActiveTime: 0,
			UUID:           UUID}
	}
	allClients = append(allClients, *c)
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
