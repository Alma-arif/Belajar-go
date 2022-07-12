package main

import (
	"belajar-go/auth"
	"belajar-go/handler"
	"belajar-go/user"
	"log"

	"github.com/gin-gonic/gin"
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
	api.POST("/avatars", userHandler.UploadAvatar)

	routes.Run()

}
