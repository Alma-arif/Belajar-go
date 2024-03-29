package main

import (
	"belajar-go/auth"
	"belajar-go/campaign"
	"belajar-go/handler"
	"belajar-go/helper"
	"belajar-go/user"
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
	campaignRepository := campaign.NewRepository(db)
	// campaign, _ := campaignRepository.FindByUserID(6)

	// for _, campaigns := range campaign {
	// 	fmt.Println(campaigns.CampaignImages[0].FileName)

	// }

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	// data, _ := campaignService.FindCampaigns(6)

	// // fmt.Println(data)
	// for _, d := range data {
	// 	fmt.Println(d)

	// }

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

	// input := campaign.CreateCampaignInput{}

	// input.Name = "ayoo bagi duit"
	// input.ShortDescription = "sokkoko"
	// input.Description = "description ayo bagi duit"
	// input.GoalAmount = 100000000
	// input.Perks = "ayo, dia, makan, aku"

	// inputUser, _ := userService.FindUserByID(6)
	// input.User = inputUser

	// _, err = campaignService.CreateCampaign(input)
	// if err != nil {
	// 	log.Fatal(err.Error())

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	routes := gin.Default()

	routes.Static("/images", "./images")
	api := routes.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	routes.Run()

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
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
