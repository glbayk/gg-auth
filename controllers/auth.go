package controllers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glbayk/gg-auth/models"
	"github.com/glbayk/gg-auth/utils"
	"github.com/go-playground/validator/v10"
)

type AuthController struct{}

type registerDto struct {
	Email    string `json:"email" binding:"required,email,min=3,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (ac AuthController) Register(ctx *gin.Context) {
	var payload registerDto

	err := ctx.BindJSON(&payload)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ErrorMessage, len(ve))

			for i, fe := range ve {
				out[i] = utils.TransformErrorMessage(fe)
			}

			ctx.JSON(http.StatusBadRequest, gin.H{"error": out})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request"})
		return
	}

	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    payload.Email,
		Password: hash,
	}

	if user.Find() == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	err = user.Create()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create user"})
		return
	}

	ctx.JSON(http.StatusOK, "User registered")
}

type loginDTO struct {
	Email    string `json:"email" binding:"required,email,min=3,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (ac AuthController) Login(ctx *gin.Context) {
	var payload loginDTO

	err := ctx.BindJSON(&payload)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ErrorMessage, len(ve))

			for i, fe := range ve {
				out[i] = utils.TransformErrorMessage(fe)
			}

			ctx.JSON(http.StatusBadRequest, gin.H{"error": out})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request"})
		return
	}

	user := models.User{
		Email: payload.Email,
	}

	err = user.Find()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	err = utils.ComparePassword(payload.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	atTtl, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	at, err := utils.TokenCreate(atTtl, user.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	rtTtl, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	rt, err := utils.TokenCreate(rtTtl, user.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": at, "refresh_token": rt})
}

func (ac AuthController) RefreshToken(ctx *gin.Context) {
	var refreshToken string

	ah := ctx.Request.Header.Get("Authorization")
	hf := strings.Fields(ah)
	if len(hf) > 0 && hf[0] == "Bearer" {
		refreshToken = hf[1]
	}

	if refreshToken == "" {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.TokenValidate(refreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	claimsEmail, err := claims.GetSubject()
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	atTtl, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	at, err := utils.TokenCreate(atTtl, claimsEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	rtTtl, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION_TIME_MINUTES"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	rt, err := utils.TokenCreate(rtTtl, claimsEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not create token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": at, "refresh_token": rt})
}

type resetPasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=72"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

func (ac AuthController) ResetPassword(ctx *gin.Context) {
	var payload resetPasswordDTO

	err := ctx.BindJSON(&payload)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ErrorMessage, len(ve))

			for i, fe := range ve {
				out[i] = utils.TransformErrorMessage(fe)
			}

			ctx.JSON(http.StatusBadRequest, gin.H{"error": out})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request"})
		return
	}

	user := ctx.MustGet("user").(models.User)
	if utils.ComparePassword(payload.OldPassword, user.Password) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Old password does not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.NewPassword)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	err = user.Find()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	user.Password = hashedPassword
	err = user.Update()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not update user password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}
