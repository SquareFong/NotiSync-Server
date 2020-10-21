package main

import (
	"encoding/json"
	"fmt"
)

type communicateStruct struct {
	UUID string
	Time string
	Type string
	Data string //base64
}

type notification struct {
	Time        string
	PackageName string
	Title       string
	Content     string
}

func strToNotification(str string) notification {
	var n notification
	err := json.Unmarshal([]byte(str), &n)
	if err != nil {
		fmt.Println("notificationParser:strToNotification:\n json ERROR", err)
	}
	return n
}

func notificationToStr(noti communicateStruct) string {
	str, err0 := json.Marshal(noti)
	if err0 != nil {
		fmt.Println("notificationParser:notificationToStr:\n json err:", err0)
	}
	return string(str)
}

func notificationsToStr(noti []notification) string {
	str, err0 := json.Marshal(noti)
	if err0 != nil {
		fmt.Println("notificationParser:notificationToStr:\n json err:", err0)
	}
	return string(str)
}

type phoneDetail struct {
	OsVersion    string
	Model        string
	Kernel       string
	Uptime       string
	Processor    string
	MemoryUsage  string
	StorageUsage string
}

func strToPhoneDetail(str string) phoneDetail {
	var item phoneDetail
	err := json.Unmarshal([]byte(str), &item)
	if err != nil {
		fmt.Println("notificationParser:strToNotification:\n json ERROR", err)
	}
	return item
}

type message struct {
	Number string
	Name   string
	Body   string
	Date   string
	Type   string
}

func strToMessage(str string) message {
	var item message
	err := json.Unmarshal([]byte(str), &item)
	if err != nil {
		fmt.Println("notificationParser:strToNotification:\n json ERROR", err)
	}
	return item
}

func strToAllMessages(str string) []message {
	var items []message
	err := json.Unmarshal([]byte(str), &items)
	if err != nil {
		fmt.Println("notificationParser:strToNotification:\n json ERROR", err)
	}
	return items
}
