package main

import (
	"crud-todo/todo/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

/*
 * Initiates the routes function and check if .env file exists
 * and port is set in .env file.
 */
func main() {
	logfile, err := os.Create("todo.log")
	if err != nil {
		fmt.Println("Log file not created %s", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("PORT not set in .env")
	}

	err = http.ListenAndServe(":"+port, routes.Init())
	if err != nil {
		log.Printf("%s", err)
	}

	log.SetOutput(logfile)
	logfile.Close()
}
