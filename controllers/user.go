package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glbayk/gg-auth/models"
)

type UserController struct{}

func (u *UserController) Profile(ctx *gin.Context) {
	user := ctx.MustGet("user").(models.User)
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
