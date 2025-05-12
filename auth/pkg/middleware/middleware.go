package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luntsev/notes-manager/auth/pkg/jwt"
)

func IsAuth(jwtServ *jwt.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authStr := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authStr, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authStr, "Bearer ")

		data, err := jwtServ.Verify(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("email", data.Email)
		ctx.Set("id", data.Id)
		ctx.Next()
	}
}
