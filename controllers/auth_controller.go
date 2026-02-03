package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"gin-api/config"
	"gin-api/helpers"
	"gin-api/services"

	"github.com/gin-gonic/gin"
)

type GoogleUser struct {
	Name          string
	Email         string
	Firstname     string
	Surname       string
	Picture       string
	VerifiedEmail bool
}

func (g *GoogleUser) BindJSON(obj interface{}) error {
	data, ok := obj.(map[string]interface{})
	if !ok {
		return nil
	}
	g.Name = helpers.GetString(data, "name")
	g.Email = helpers.GetString(data, "email")
	g.Firstname = helpers.GetString(data, "given_name")
	g.Surname = helpers.GetString(data, "family_name")
	g.Picture = helpers.GetString(data, "picture")
	g.VerifiedEmail = helpers.GetBool(data, "verified_email")
	return nil
}

func GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}

	// 1. exchange token
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed"})
		return
	}

	// 2. get user info from google
	client := config.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var rawUser map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decode google user failed"})
		return
	}

	user, err := services.FindOrCreateGoogleUser(rawUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "login failed"})
		return
	}

	jwtToken, err := helpers.GenerateJWT(user.ID.Hex(), user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"data":  user,
		"code":  http.StatusOK,
	})

}
