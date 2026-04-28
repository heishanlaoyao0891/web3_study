package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID            uint
	Username      string
	Level         int
	Exp           int
	Coins         int
	CheckinDays   int
	LastCheckinAt *string
}

type Checkin struct {
	ID          uint
	UserID      uint
	CheckinDate string
	ExpGained   int
	CoinsGained int
}

func main() {
	Db, err := gorm.Open(mysql.Open("root:Woshizhu?@tcp(124.223.6.26:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	var users []User
	Db.Table("users").Select("id, username, level, exp, coins, checkin_days, last_checkin_at").Scan(&users)

	fmt.Println("用户数据：")
	for _, u := range users {
		lastCheckin := "NULL"
		if u.LastCheckinAt != nil {
			lastCheckin = *u.LastCheckinAt
		}
		fmt.Printf("ID:%d 用户:%s Lv:%d 经验:%d 金币:%d 打卡天数:%d 最后打卡:%s\n",
			u.ID, u.Username, u.Level, u.Exp, u.Coins, u.CheckinDays, lastCheckin)
	}

	fmt.Println("\n打卡记录：")
	var checkins []Checkin
	Db.Table("checkins").Order("id desc").Limit(10).Scan(&checkins)
	for _, c := range checkins {
		fmt.Printf("ID:%d 用户ID:%d 日期:%s 经验:%d 金币:%d\n",
			c.ID, c.UserID, c.CheckinDate, c.ExpGained, c.CoinsGained)
	}
}
