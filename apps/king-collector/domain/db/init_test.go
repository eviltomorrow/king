package db

import (
	"log"
	"time"

	"github.com/eviltomorrow/king/lib/db/mongodb"
)

func init() {
	_, err := mongodb.InitMongoDB(&mongodb.Config{
		DSN:            "mongodb://admin:admin123@mongo:27017/transaction_db",
		MinOpen:        3,
		MaxOpen:        10,
		ConnectTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
}
