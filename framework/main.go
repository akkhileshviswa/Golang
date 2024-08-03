package main

import (
	"Golang/framework/handler"
	"Golang/framework/middleware"
	"Golang/framework/model"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")

	r := SetupRouter()

	log.Fatal(r.Run(":" + port))
}

func DbInit() *gorm.DB {
	db, err := model.Setup()
	if err != nil {
		log.Println("Problem setting up database")
	}

	return db
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	db := DbInit()

	server := handler.NewServer(db)

	router := r.Group("/api")
	router.POST("/register", server.Register)
	router.POST("/login", server.Login)

	authorized := r.Group("/api/admin")
	authorized.Use(middleware.JwtAuthMiddleware())
	authorized.GET("/groceries", server.GetGroceries)
	authorized.POST("/grocery", server.PostGrocery)

	return r
}
