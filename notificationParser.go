package main

import (
	"encoding/json"
	"fmt"
)

type notification struct {
	UUID string
	Time string
	Data []notificationData
}

type notificationData struct {
	Time string
	PackageName string
	Title string
	Content string
}

func strToNotification(str string) notification {
	var n notification
	err := json.Unmarshal([]byte(str), &n)
	if err != nil {
		fmt.Println("notificationParser:strToNotification:\n json ERROR", err)
	}
	return n
}

func notificationToStr(noti notification) string  {
	str, err0 := json.Marshal(noti)
	if err0 != nil {
		fmt.Println("notificationParser:notificationToStr:\n json err:", err0)
	}
	return string(str)
}

func notificationsToStr(noti []notificationData) string  {
	str, err0 := json.Marshal(noti)
	if err0 != nil {
		fmt.Println("notificationParser:notificationToStr:\n json err:", err0)
	}
	return string(str)
}