package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/adsense/v1.4"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	clientId := os.Getenv("GOOGLE_KEY")
	clientSecret := os.Getenv("GOOGLE_SECRET")
	redirectUri := os.Getenv("GOOGLE_REDIRECT_URI")

	if len(clientId) == 0 {
		log.Fatal("GOOGLE_CLIENT_ID is empty.")
	}

	if len(clientSecret) == 0 {
		log.Fatal("GOOGLE_CLIENT_SECRET is empty.")
	}

	if len(redirectUri) == 0 {
		log.Fatal("GOOGLE_REDIRECT_URI is empty.")
	}

	oauth := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectUri,
		Scopes: []string{
			"https://www.googleapis.com/auth/adsense.readonly",
		},
	}

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		url := oauth.AuthCodeURL("state", oauth2.AccessTypeOffline)
		ctx.Redirect(http.StatusMovedPermanently, url)
	})

	router.GET("/auth", func(ctx *gin.Context) {
		code := ctx.Query("code")

		token, err := oauth.Exchange(oauth2.NoContext, code)
		if err != nil {
			ctx.JSON(500, err)
			return
		}

		client := oauth.Client(oauth2.NoContext, token)
		service, err := adsense.New(client)

		call := service.Accounts.List()
		resp, err := call.Do()
		if err != nil {
			ctx.JSON(500, err)
			return
		}

		ctx.JSON(200, resp)
	})

	router.Run(":8080")
}
