package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"gin-api/config"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "code not found"})
		return
	}

	token, err := config.GoogleOAuthConfig.Exchange(
		context.Background(),
		code,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "token exchange failed"})
		return
	}

	client := config.GoogleOAuthConfig.Client(
		context.Background(),
		token,
	)

	resp, err := client.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo",
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)

	c.JSON(200, gin.H{"data": gin.H{
		"email":          user["email"],
		"surname":        user["family_name"],
		"firstname":      user["given_name"],
		"id":             user["id"],
		"full_name":      user["name"],
		"picture":        user["picture"],
		"verified_email": user["verified_email"],
	}, "code": 200})
}
