package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var dataSourceName string

func openDB() *sql.DB {
	db, err := sql.Open(
		"mysql",
		dataSourceName)
	checkErr(err)
	return db
}

func readDBConfig(fileName string) string {
	fl, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer fl.Close()
	//读文件
	var str string
	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if 0 == n {
			break
		}
		str += string(buf)
	}
	str = strings.Replace(str, "\x00", "", -1)
	//转Json
	type database struct {
		UserName string
		Password string
		DBName   string
	}
	var dbCfg database
	err0 := json.Unmarshal([]byte(str), &dbCfg)
	if err0 != nil {
		fmt.Println("notificationParser:strToNotification:\n json ERROR", err0)
	} else {
		println(dbCfg.UserName + " " + dbCfg.Password + " " + dbCfg.DBName)
	}
	r := fmt.Sprintf("%s:%s@/%s?charset=utf8", dbCfg.UserName, dbCfg.Password, dbCfg.DBName)
	println("readDBConfig: " + r)
	dataSourceName = r
	return r
}

func createUsersTable() {
	db := openDB()
	s := "create table if not exists Users (id INTEGER key AUTO_INCREMENT, uuid TEXT);"
	smt, err := db.Prepare(s)
	checkErr(err)
	smt.Exec()
}

func createNotiTable(id int) {
	db := openDB()
	s := "create table if not exists " + "table" + string(id) +
		" (time INTEGER, PackageName TEXT, Title TEXT, Content TEXT)"
	smt, err := db.Prepare(s)
	checkErr(err)
	smt.Exec()
}

func addUser(uuid string) {
	db := openDB()
	s := fmt.Sprintf("insert Users(uuid) values('%s')", uuid)
	smt, err := db.Prepare(s)
	checkErr(err)
	smt.Exec()
}

func getUserCounts() int {
	db := openDB()
	s := "SELECT COUNT(*) FROM Users"
	rows, err := db.Query(s)
	defer rows.Close()
	checkErr(err)
	n := 0
	for rows.Next() {
		rows.Scan(&n)
	}
	return n
}

// 加入一个通知, 按照uuid
func insertNotificationByUUID(UUID string, data notificationData) {
	db := openDB()
	tableName := getTableName(UUID)
	var Time int
	Time, err := strconv.Atoi(data.Time)
	checkErr(err)
	s := fmt.Sprintf(
		"insert %s(time, PackageName, Title, Content) values(%d, '%s', '%s', '%s')",
		tableName, Time, data.PackageName, data.Title, data.Content)
	smt, err := db.Prepare(s)
	checkErr(err)
	smt.Exec()
}

// 获取通知
func getNotification(UUID string, lastUpdate string) []notificationData {
	db := openDB()
	tableName := getTableName(UUID)
	s := fmt.Sprintf("select * from %s where time > %s", tableName, lastUpdate)
	rows, err := db.Query(s)
	checkErr(err)
	defer rows.Close()
	var datas []notificationData
	for rows.Next() {
		var data notificationData
		rows.Scan(&data.Time, &data.PackageName, &data.Title, &data.Content)
		datas = append(datas, data)
	}
	return datas
}

func getTableName(uuid string) string {
	db := openDB()
	s := "select id from Users where uuid='" + uuid + "'"
	rows, err := db.Query(s)
	checkErr(err)
	defer rows.Close()
	str := ""
	for rows.Next() {
		var id int
		rows.Scan(&id)
		str = "table" + strconv.Itoa(id)
	}
	return str
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
