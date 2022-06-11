package main

import (
	db1 "github.com/hktalent/goSqlite_gorm/pkg/db"
	"gorm.io/gorm"
	"log"
	"time"
)

// 记录本地解析
type Result struct {
	gorm.Model
	Dns    string `json:"domain" gorm:"unique_index"`
	Ips    []Ips  `json:"ips" gorm:"many2many:result_ips"`
	Date   string `json:"date"`
	SaveEs bool   `json:"saveEs"`
}

type Ips struct {
	gorm.Model
	Ip string `json:"ip" gorm:"unique_index"`
}

func getDate() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 10:04:05")
}
func main() {
	log.Println(getDate())
	xx0 := db1.GetDb(&Ips{}, "/Users/51pwn/MyWork/DNS_Server/db/mydbfile")
	xx0.AutoMigrate(&Ips{}, &Result{})
	db1.GetTableName(&Result{})
	var r Result = Result{Ips: []Ips{Ips{Ip: "11.2.33.4"}}, Dns: "xxx.com", Date: getDate()}
	if 1 == db1.Create[Result](&r) {
		log.Println("db1.Create ok ")
	}
	var rst []Result = db1.GetSubQueryLists(Result{}, "Ips",
		[]Result{}, 10, 0, "dns = ?", "xxx.com")
	log.Println("rst len = ", len(rst))
	log.Println(rst[0])
	var r3 Result
	r1 := db1.GetOne[Result](&r3, "dns=?", "xxx.com")
	log.Println(r1)
	log.Println(r3.Ips)
}
