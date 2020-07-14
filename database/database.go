package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go-rest/app"
	"log"
	"os"
)

var (
	DBConn *gorm.DB
)

func InitDatabase() {
	var err error

	dbFile, dbLocExist := os.LookupEnv(app.EnvDatabaseFile)
	if !dbLocExist {
		dbFile = app.DefaultDBFile
	}

	DBConn, err = gorm.Open(app.DBDialect, dbFile)
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database successfully connected")
}

func CloseDatabase() {
	log.Fatal(DBConn.Close())
}
