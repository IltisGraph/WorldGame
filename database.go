package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func InitDB(host, user, password string) {
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=worldgamedb port=5432 sslmode=disable TimeZone=Europe/Berlin"
	a, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	DB = a

	AutoMigrate()
}