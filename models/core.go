package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	fmt.Println("init db")
	username := "root"
	password := "123456"
	host := "127.0.0.1"
	port := 3306
	dbname := "keeper"
	timeout := "10s"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, dbname, timeout)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err := DB.AutoMigrate(&Account{})
	if err != nil {
		fmt.Println("account create table err:", err)
		panic(err)
		return
	}
}
