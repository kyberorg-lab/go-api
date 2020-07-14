package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go-rest/app"
	"log"
)

var (
	DBConn *gorm.DB
)

func InitDatabase() {
	var err error
	DBConn, err = gorm.Open(app.DBDialect, app.DBFile)
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database successfully connected")
}

func CloseDatabase() {
	log.Fatal(DBConn.Close())
}
