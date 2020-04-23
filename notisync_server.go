package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func parse(writer http.ResponseWriter, request *http.Request){
	request.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	fmt.Println("Form: ", request.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("Method: ", request.Method)
	fmt.Println("Path: ", request.URL.Path)
	if request.Method == "POST" {
		if request.URL.Path != "/send" {
			fmt.Println("Path and Method did not match")
			return
		}
		//处理POST内容
		result, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal("parse: ", err)
		}
		str := bytes.NewBuffer(result).String()
		fmt.Println(str)
		var noti notification
		err0 := json.Unmarshal([]byte(str), &noti)
		if err0 != nil {
			fmt.Println("notificationParser:strToNotification:\n json ERROR", err0)
		}
		//测试用打印
		fmt.Println(noti)

		if len(noti.Data) != 0 {
			for _, n := range noti.Data {
				insertNotificationByUUID(noti.UUID, n)
			}
		}else {
			fmt.Println("No notification received!")
		}
		fmt.Fprintf(writer, "200")
	} else if request.Method == "GET" {
		if request.URL.Path != "/get" {
			fmt.Println("Path and Method did not match")
			return
		}
		for k, v := range request.Form {
			fmt.Println("key", k)
			fmt.Println("val: ", strings.Join(v, ""))
		}
		uuid := request.Form.Get("UUID")
		lastUpdate := request.Form.Get("Time")
		if len(uuid) == 0 {
			fmt.Println("UUID is Null")
		}else if len(lastUpdate) == 0 {
			fmt.Println("Time is Null")
		} else {
			noti := getNotification(uuid, lastUpdate)
			str := notificationsToStr(noti)
			fmt.Println(str)
			n, err := fmt.Fprintf(writer, str)
			fmt.Println(n)
			fmt.Println(err)
		}
	} else {
		println("Got Other Methods")
	}

}

func main() {
	createUsersTable()

	http.HandleFunc("/", parse) //设置访问的路由
	err := http.ListenAndServe(":9000", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}