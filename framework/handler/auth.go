package handler

import (
	"Golang/framework/model"
	"Golang/framework/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (s *Server) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{Username: input.Username, Password: input.Password}
	user.HashPassword()

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "message": "User created"})
}

func (s *Server) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{Username: input.Username, Password: input.Password}

	token, err := s.LoginCheck(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The username or password is not correct"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) CurrentUser(c *gin.Context) {
	err := utils.ValidateToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GetToken(c)
	if err != nil {
		log.Println(err)
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	user, err := model.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": user})
}

func (s *Server) LoginCheck(username, password string) (string, error) {
	var err error

	user := model.User{}

	if err = s.db.Model(model.User{}).Where("username=?", username).Take(&user).Error; err != nil {
		return "", err
	}

	err = model.VerifyPassword(password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (s *Server) GetUserById(c *gin.Context) {
	var user model.User

	if err := s.db.Model(model.User{}).Where("ID=?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": user})
}
