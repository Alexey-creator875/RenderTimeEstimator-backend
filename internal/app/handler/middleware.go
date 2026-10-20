package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"

	"RenderTimeEstimator/internal/app/ds"
	"RenderTimeEstimator/internal/app/role"
)

const jwtPrefix = "Bearer "

func (h *Handler) WithAuthCheck(assignedRoles ...role.Role) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		jwtStr = jwtStr[len(jwtPrefix):]

		// проверяем jwt в блеклисте redis
		err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil { // ключ найден, значит токен в блеклисте
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !errors.Is(err, goredis.Nil) { // не "ключа нет", а настоящая ошибка redis
			logrus.Error(err)
			gCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(h.JWT.Token), nil
		})
		if err != nil || !token.Valid {
			logrus.Error(err)
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		myClaims := token.Claims.(*ds.JWTClaims)

		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				gCtx.Set("userID", myClaims.UserID)
				gCtx.Set("userUUID", myClaims.UserUUID)
				gCtx.Set("role", myClaims.Role)
				gCtx.Next()
				return
			}
		}

		logrus.Printf("role %v is not assigned in %v", myClaims.Role, assignedRoles)
		gCtx.AbortWithStatus(http.StatusForbidden)
	}
}

func (h *Handler) WithOptionalAuth(gCtx *gin.Context) {
	jwtStr := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		gCtx.Next()
		return
	}
	jwtStr = jwtStr[len(jwtPrefix):]

	// токен из блеклиста считаем отсутствующим
	if err := h.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr); err == nil {
		gCtx.Next()
		return
	}

	token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.JWT.Token), nil
	})
	if err == nil && token.Valid {
		claims := token.Claims.(*ds.JWTClaims)
		gCtx.Set("userID", claims.UserID)
		gCtx.Set("role", claims.Role)
	}

	gCtx.Next()
}