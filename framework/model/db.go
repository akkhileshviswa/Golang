package model

import (
	"fmt"

	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
 * This function is used to inialize the database connection
 *
 * @return gorm.DB
 */
func Setup() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	if user == "" {
		log.Fatal("DB_USER not set in .env")
	}

	pass := os.Getenv("DB_PASS")
	if pass == "" {
		log.Fatal("DB_PASS not set in .env")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Fatal("DB_HOST not set in .env")
	}

	dbUrl := fmt.Sprintf("%s:%s@(%s:3306)/%s?charset=utf8&parseTime=True", user, pass, host, "ormdemo")
	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.AutoMigrate(&User{}); err != nil {
		log.Println(err)
	}

	if err = db.AutoMigrate(&Grocery{}); err != nil {
		log.Println(err)
	}

	return db, err
}
