package handler

import (
	"RenderTimeEstimator/internal/app/ds"
	"RenderTimeEstimator/internal/app/role"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	login    = "login"
	password = "password"
)

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.Login != login || req.Password != password {
		c.AbortWithStatus(http.StatusForbidden)
		logrus.Println("sdklf")
		return
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(h.JWT.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "render-time-estimator",
		},
		UserUUID: uuid.New(), // тестовый uuid
	})

	str, err := token.SignedString([]byte(h.JWT.Token))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, loginResp{
		ExpiresIn:   h.JWT.ExpiresIn,
		AccessToken: str,
		TokenType:   "Bearer",
	})
}

type registerReq struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type registerResp struct {
	Ok bool `json:"ok"`
}

func (h *Handler) Register(gCtx *gin.Context) {
	req := &registerReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.Password == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("pass is empty"))
		return
	}

	if req.Login == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}

	err = h.Repository.Register(&ds.User{
		UUID: uuid.New(),
		Role: role.Guest,
		Login: req.Login,
		Password: generateHashString(req.Password), // пароли делаем в хешированном виде и далее будем сравнивать хеши, чтобы их не угнали с базой вместе
	})
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	gCtx.JSON(http.StatusOK, &registerResp{
		Ok: true,
	})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *Handler) SignUp(ctx *gin.Context) {
	newUser := ds.User{
		Login:		ctx.Request.FormValue("login"),
		Password:	ctx.Request.FormValue("password"),
	}

	err := h.Repository.AddUser(&newUser)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user": newUser,
	})
}

func (h *Handler) SignIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "авторизация",
	})
}

func (h *Handler) SignOut(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "деавторизация",
	})
}