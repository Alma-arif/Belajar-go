package main

import (
	"belajar-go/auth"
	"belajar-go/handler"
	"belajar-go/helper"
	"belajar-go/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/ayo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	// fmt.Println(authService.GenerateToken(10))

	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2fQ._ZqlWXLmUva8BDYev0XWcoN0xYVMpvhlsZUypT_96oY")

	if err != nil {
		fmt.Println("erroe")
		fmt.Println("erroe")

	}

	if token.Valid {
		fmt.Println("valid")
		fmt.Println("valid")
	} else {
		fmt.Println("invalid")
		fmt.Println("invalid")
	}

	// userService.SaveAvatar(1, "images/1-profile.png")
	// cek service login
	// input := user.LoginInput{
	// 	Email:    "aku@mail.com",
	// 	Password: "passworda",
	// }

	// user, err := userService.Login(input)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println(user.Email)

	// cek repository by email
	// userByEmail, err := userRepository.FindByEmail("@mail.com")
	// if err != nil {
	// 	fmt.Println(er r.Error())
	// }

	// if userByEmail.ID == 0 {
	// 	fmt.Println("user tidak di temukan")
	// } else {
	// 	fmt.Println(userByEmail.Email)
	// }

	userHandler := handler.NewUserHandler(userService, authService)

	routes := gin.Default()

	api := routes.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddlemare(authService, userService), userHandler.UploadAvatar)

	routes.Run()

}

func authMiddlemare(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		// id dari jwt untuk  di cek di servic user bedasarkan id
		user, err := userService.FindUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}
}
