package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glbayk/gg-auth/models"
	"github.com/glbayk/gg-auth/utils"
)

func Authenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		ah := ctx.Request.Header.Get("Authorization")
		hf := strings.Fields(ah)
		if len(hf) > 0 && hf[0] == "Bearer" {
			accessToken = hf[1]
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		claims, err := utils.TokenValidate(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		claimsEmail, err := claims.GetSubject()
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		user := models.User{Email: claimsEmail}
		user.Find()

		ctx.Set("user", user)
		ctx.Next()
	}
}
