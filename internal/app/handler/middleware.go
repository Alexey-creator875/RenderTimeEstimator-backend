package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"

	"RenderTimeEstimator/internal/app/ds"
)

const jwtPrefix = "Bearer "

func (h *Handler) WithAuthCheck(c *gin.Context) {
	s := c.GetHeader("Authorization")
	if !strings.HasPrefix(s, jwtPrefix) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	s = s[len(jwtPrefix):]

	token, err := jwt.ParseWithClaims(s, &ds.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.JWT.Token), nil
	})
	if err != nil || !token.Valid {
		logrus.Error(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	claims := token.Claims.(*ds.JWTClaims)
	c.Set("userUUID", claims.UserUUID)
	c.Next()
}