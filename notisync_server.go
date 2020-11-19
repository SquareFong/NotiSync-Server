package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func parse(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()                     //解析url传递的参数，对于POST则解析响应包的主体（request body）
	fmt.Println("Method: ", request.Method) //这些信息是输出到服务器端的打印信息
	fmt.Println("Path: ", request.URL.Path)
	fmt.Println("Form: ", request.Form)
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
		var commStruct communicateStruct
		err0 := json.Unmarshal([]byte(str), &commStruct)
		if err0 != nil {
			fmt.Println("notificationParser:strToNotification:\n json ERROR", err0)
		}
		//测试用打印
		fmt.Println(commStruct)

		if commStruct.Type == "Notifications" {
			if len(commStruct.Data) != 0 {
				decodeByte, _ := base64.StdEncoding.DecodeString(commStruct.Data)
				n := strToNotification(string(decodeByte))

				fmt.Println(n)
				//入数据库
				if n != nil {
					for _, item := range n {
						insertNotificationByUUID(commStruct.UUID, item)
					}
				}
			} else {
				fmt.Println("No communicateStruct received!")
			}
		} else if commStruct.Type == "Detail" {
			if len(commStruct.Data) != 0 {
				decodeByte, _ := base64.StdEncoding.DecodeString(commStruct.Data)
				detail := strToPhoneDetail(string(decodeByte))
				fmt.Println(detail)
				setDetail(commStruct.UUID, detail)
			} else {
				fmt.Println("No communicateStruct received!")
			}
		} else if commStruct.Type == "Messages" {
			if len(commStruct.Data) != 0 {
				decodeByte, _ := base64.StdEncoding.DecodeString(commStruct.Data)
				msg := strToMessages(string(decodeByte))
				fmt.Println(msg)
				setMessages(commStruct.UUID, msg)
			} else {
				fmt.Println("No communicateStruct received!")
			}

		} else if commStruct.Type == "AllMessages" {
			if len(commStruct.Data) != 0 {
				decodeByte, _ := base64.StdEncoding.DecodeString(commStruct.Data)
				msg := strToMessages(string(decodeByte))
				fmt.Println(msg)
				setAllMessage(commStruct.UUID, msg)
			} else {
				fmt.Println("No communicateStruct received!")
			}
		} else if commStruct.Type == "newSMS" {
			if len(commStruct.Data) != 0 {
				decodeByte, _ := base64.StdEncoding.DecodeString(commStruct.Data)
				msg := strToMessage(string(decodeByte))
				fmt.Println(msg)
				setNewSMS(commStruct.UUID, msg)
			} else {
				fmt.Println("No communicateStruct received!")
			}
		} else {

		}

		fmt.Fprintf(writer, "200")
	} else if request.Method == "GET" {
		//Get request
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
		cmdType := request.Form.Get("Type")
		if len(uuid) == 0 {
			fmt.Println("UUID is Null")
		} else if len(lastUpdate) == 0 {
			fmt.Println("Time is Null")
		} else {
			if cmdType == "Notifications" {
				noti := getNotification(uuid, lastUpdate)
				data := notificationsToStr(noti)
				fmt.Println(data)
				comStructStr := packageToCommStr(uuid, strconv.FormatInt(time.Now().Unix(), 10),
					"Notifications", data)
				fmt.Println(comStructStr)
				_, err := fmt.Fprintf(writer, comStructStr)
				if err != nil {
					fmt.Println(err)
				}
			} else if cmdType == "Detail" {
				//TODO
				detail := getDetail(uuid)
				data := detailToStr(*detail)
				comStructStr := packageToCommStr(uuid, strconv.FormatInt(time.Now().Unix(), 10),
					"Detail", data)
				_, _ = fmt.Fprintf(writer, comStructStr)
			} else if cmdType == "Messages" {
				msm := getMessages(uuid)
				data := messagesToStr(*msm)
				comStructStr := packageToCommStr(uuid, strconv.FormatInt(time.Now().Unix(), 10),
					"Messages", data)
				_, _ = fmt.Fprintf(writer, comStructStr)
			}  else if cmdType == "AllMessages" {
				msm := getMessages(uuid)
				data := messagesToStr(*msm)
				comStructStr := packageToCommStr(uuid, strconv.FormatInt(time.Now().Unix(), 10),
					"AllMessages", data)
				_, _ = fmt.Fprintf(writer, comStructStr)
			} else if cmdType == "Command" {
				var c *client
				var t string
				var data = ""
				if findClientByUUID(uuid, &c) {
					if len(c.newSMS) != 0 {
						data = messagesToStr(*getNewMessages(uuid))
						c.newSMS = nil
					} else if c.needAll {
						t = "all"
						c.needAll = false
					} else if time.Now().Unix()-c.lastActiveTime > 100 {
						//TODO 删结构体
						t = "dead"
					} else {
						t = "active"
					}
				} else {
					t = "dead"
				}
				if data != "" {
					data = base64.StdEncoding.EncodeToString([]byte(data))
				}
				comStructStr := packageToCommStr(uuid, strconv.FormatInt(time.Now().Unix(), 10),
					t, data)
				_, _ = fmt.Fprintf(writer, comStructStr)
			} else {
				var str string
				n, err := fmt.Fprintf(writer, str)
				fmt.Println(n)
				fmt.Println(err)
			}
		}
	} else {
		println("Got Other Methods")
	}

}

func main() {
	cfgFilePtr := flag.String("c", "", "configure file position")

	flag.Parse()

	if len(*cfgFilePtr) == 0 {
		*cfgFilePtr = "/etc/notisync/config.json"
	}

	readDBConfig(*cfgFilePtr)

	createUsersTable()

	http.HandleFunc("/", parse)              //设置访问的路由
	err := http.ListenAndServe(":9000", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
