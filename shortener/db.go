package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type URL struct {
	gorm.Model
	Shortened string `gorm:"unique_index"`
	URL       string
}

var DB *gorm.DB

func openDB(dbname string) {
	db, err := gorm.Open("sqlite3", dbname)
	if err != nil {
		panic("failed to open database")
	}
	DB = db

	db.AutoMigrate(&URL{})
}

func closeDB() {
	DB.Close()
}

func shortenURL(toShorten string) string {
	url := URL{URL: toShorten}
	fmt.Println("shorten ", toShorten)

	// check if url is already shortened
	DB.Where(&URL{URL: toShorten}).First(&url)
	if url.Shortened != "" {
		return url.Shortened
	}

	for {
		var u URL
		url.Shortened = generateShortID(5)
		DB.Where((&URL{Shortened: url.Shortened})).First(&u)
		if u.URL == "" {
			break
		}
	}

	DB.Create(&url)
	return url.Shortened
}

func getURL(ID string) (string, error) {
	var url URL
	DB.Where((&URL{Shortened: ID})).First(&url)
	if url.URL != "" {
		return url.URL, nil
	}
	return "", errors.New("Can't find URL")
}

func generateShortID(length int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"
	const alphlen = len(alphabet)
	var r strings.Builder
	for i := 0; i < length; i++ {
		r.WriteByte(alphabet[rand.Intn(alphlen)])
	}
	return r.String()
}
