package middleware

import (
	"strings"

	"github.com/JordanMarcelino/go-gin-starter/internal/constant"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/httperror"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/utils/jwtutils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtUtil jwtutils.JwtUtil
}

func NewAuthMiddleware(jwtUtil jwtutils.JwtUtil) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtil: jwtUtil,
	}
}

func (m *AuthMiddleware) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := m.parseAccessToken(ctx)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		claims, err := m.jwtUtil.Parse(accessToken)
		if err != nil {
			ctx.Error(httperror.NewUnauthorizedError())
			ctx.Abort()
			return
		}

		ctx.Set(constant.CTX_USER_ID, claims.UserID)
		ctx.Next()
	}
}

func (m *AuthMiddleware) parseAccessToken(ctx *gin.Context) (string, error) {
	accessToken := ctx.GetHeader("Authorization")
	if accessToken == "" || len(accessToken) == 0 {
		return "", httperror.NewUnauthorizedError()
	}

	splitToken := strings.Split(accessToken, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		return "", httperror.NewUnauthorizedError()
	}
	return splitToken[1], nil
}
