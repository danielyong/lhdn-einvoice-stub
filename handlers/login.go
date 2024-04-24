package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type lhdnLogin struct {
	ClientId    string `gorm:"primaryKey"`
	AccessToken string
	TokenExpiry int64 `gorm:"autoCreateTime"`
}

type lhdnAuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func ValidateAccessToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	if token == "" {
		ctx.JSON(401, &lhdnError{
			Code:    "Unauthorized",
			Message: "Authorization token not provided",
		})
    ctx.Abort()
	}

	ctx.Next()
}

func IntermediaryLogin(ctx *gin.Context) {
	clientId := ctx.GetHeader("client_id")

	var existingLogin lhdnLogin
	response := lhdnAuthenticationResponse{
		TokenType: "Bearer",
		Scope:     "InvoicingApi",
	}

	var hasExistingLogin bool
	if result := db.First(&existingLogin, "`client_id`", clientId); result.RowsAffected > 0 {
		if expiryTime := time.Unix(existingLogin.TokenExpiry, 0); expiryTime.After(time.Now()) {
			response.AccessToken = existingLogin.AccessToken
			response.ExpiresIn = existingLogin.TokenExpiry - time.Now().Unix()

			ctx.JSON(200, response)
			return
		}
		hasExistingLogin = true
	}

	if response.AccessToken == "" {
		existingLogin.TokenExpiry = time.Now().Add(5 * time.Second).Unix()
		response.ExpiresIn = 60
		token := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
			"info": "This is a dummy",
			"exp":  existingLogin.TokenExpiry,
		})

		signedToken, _ := token.SignedString([]byte("something"))

		response.AccessToken = signedToken
		existingLogin.AccessToken = response.AccessToken
	}

	existingLogin.ClientId = clientId

	if hasExistingLogin {
		db.Save(&existingLogin)
	} else {
		db.Create(&existingLogin)
	}

	ctx.JSON(200, response)
}
