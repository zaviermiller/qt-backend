package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"github.com/joho/godotenv"
	"fmt"
	"time"
)

// Custom Base Model with appropriately named JSON fields
type QTModel struct {
	ID uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

var db *gorm.DB

func init() {

	fmt.Println("[*] Initializing database...")

	e := godotenv.Load()

	if e != nil {
		fmt.Print("Loaded env variables")
	}
	
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open("postgres", dbUri)

	if err != nil {
		fmt.Print(err)
	}

	db = conn
	fmt.Println("[*] Database online")
	db.Debug().AutoMigrate(&User{}, &Test{}, &Section{})

}

func GetDB() *gorm.DB {
	return db
}