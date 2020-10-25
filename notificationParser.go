package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

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

func notificationsToStr(noti []notification) string {
	str, err0 := json.Marshal(noti)
	if err0 != nil {
		fmt.Println("notificationParser:notificationsToStr:\n json err:", err0)
	}
	return string(str)
}

type communicateStruct struct {
	UUID string
	Time string
	Type string
	Data string //base64
}

//Data is raw, not base64
func packageToCommStr(UUID string, Time string, Type string, Data string) string {
	com := communicateStruct{
		UUID: UUID,
		Time: Time,
		Type: Type,
		Data: base64.StdEncoding.EncodeToString([]byte(Data))}

	str, err0 := json.Marshal(com)
	if err0 != nil {
		fmt.Println("notificationParser: communicateStructToStr:\n json err:", err0)
	}
	return string(str)

}

func communicateStructToStr(n communicateStruct) string {
	str, err0 := json.Marshal(n)
	if err0 != nil {
		fmt.Println("notificationParser: communicateStructToStr:\n json err:", err0)
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
		fmt.Println("notificationParser.go: strToPhoneDetail:\n json ERROR", err)
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
		fmt.Println("notificationParser.go: strToMessage:\n json ERROR", err)
	}
	return item
}

func strToAllMessages(str string) []message {
	var items []message
	err := json.Unmarshal([]byte(str), &items)
	if err != nil {
		fmt.Println("notificationParser.go: strToAllMessages:\n json ERROR", err)
	}
	return items
}

func allMessagesToStr(messages []message) string {
	str, err0 := json.Marshal(messages)
	if err0 != nil {
		fmt.Println("notificationParser.go: AllMessagesToStr:\n json err:", err0)
	}
	return string(str)
}

func detailToStr(detail phoneDetail) string {
	str, err0 := json.Marshal(detail)
	if err0 != nil {
		fmt.Println("notificationParser.go: detailToStr:\n json err:", err0)
	}
	return string(str)
}
