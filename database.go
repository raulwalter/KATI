package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTRGRES_USER"),
		os.Getenv("POSTRGRES_PWD"),
		os.Getenv("POSTRGRES_HOST"),
		os.Getenv("POSTRGRES_DB")),
	)

	if err != nil {
		log.Fatal(err)
	}

	db = conn

	db.AutoMigrate(&DiaryEntry{})
}

func saveDiaryEntry(entry *DiaryEntry) error {

	if entry.ID != 0 {
		return errors.New("Update not allowed")
	}

	if entry.UserName == "" {
		return errors.New("Missing user")
	}

	return db.Create(entry).Error
}

func (u *User) getDiaryEntries() ([]DiaryEntry, error) {
	entries := []DiaryEntry{}
	err := db.Where(&DiaryEntry{UserName: u.Username}).Find(&entries).Error
	return entries, err
}
