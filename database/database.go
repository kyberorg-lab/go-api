package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kyberorg/go-api/database/model"
	"github.com/kyberorg/go-api/global"
	"log"
	"os"
)

var (
	DBConn *gorm.DB
)

func InitDatabase() {
	var err error

	dbFile, dbLocExist := os.LookupEnv(global.EnvDatabaseFile)
	if !dbLocExist {
		dbFile = global.DefaultDBFile
	}

	DBConn, err = gorm.Open(global.DBDialect, dbFile)
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database successfully connected. Database location:", dbFile)

	//auto migrate
	DBConn.AutoMigrate(&model.Scope{}, &model.User{}, &model.Token{})
	fmt.Println("Database migrations are executed")
}

func CloseDatabase() {
	log.Fatal(DBConn.Close())
}
