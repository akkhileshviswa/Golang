package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

/*
 * This function is used to inialize the database connection
 *
 * @return sql.DB
 */
func Database() *sql.DB {
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

	credentials := fmt.Sprintf("%s:%s@(%s:3306)/?charset=utf8&parseTime=True", user, pass, host)

	database, err := sql.Open("mysql", credentials)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Database Connection Successful")
	}

	_, err = database.Exec(`CREATE DATABASE IF NOT EXISTS gotodo`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Exec(`USE gotodo`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
		    id INT AUTO_INCREMENT,
		    item TEXT NOT NULL,
		    completed BOOLEAN DEFAULT FALSE,
		    PRIMARY KEY (id)
		);
	`)

	if err != nil {
		fmt.Println(err)
	}

	return database
}
