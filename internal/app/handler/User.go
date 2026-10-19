package handler

import (
	"RenderTimeEstimator/internal/app/ds"
	"RenderTimeEstimator/internal/app/role"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	ExpiresIn   time.Duration `json:"expires_in" swaggertype:"integer"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}


// Login godoc
// @Summary      Аутентификация
// @Description  Проверяет логин и пароль и возвращает JWT. Токен передаётся в заголовке Authorization: Bearer <токен>.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      loginReq  true  "Логин и пароль"
// @Success      200   {object}  loginResp
// @Failure      400   "Некорректное тело запроса"
// @Failure      403   "Неверный логин или пароль"
// @Failure      500   "Ошибка сервера"
// @Router       /api/sign-in [post]
func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByLogin(req.Login)

	if err != nil || generateHashString(req.Password) != user.Password {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(h.JWT.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "render-time-estimator",
		},
		UserUUID: uuid.New(),
		Role:     user.Role,
	})

	strToken, err := token.SignedString([]byte(h.JWT.Token))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, loginResp{
		ExpiresIn:   h.JWT.ExpiresIn,
		AccessToken: strToken,
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

// Register godoc
// @Summary      Регистрация
// @Description  Создаёт нового пользователя с ролью Guest. Пароль хранится в виде хеша.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      registerReq  true  "Логин и пароль"
// @Success      200   {object}  registerResp
// @Failure      400   "Пустой логин или пароль"
// @Failure      500   "Ошибка сервера"
// @Router       /api/sign-up [post]
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

// Logout godoc
// @Summary      Выход
// @Description  Добавляет текущий JWT в блеклист Redis до истечения срока его жизни. После этого токен перестаёт приниматься.
// @Tags         Auth
// @Security     ApiKeyAuth
// @Success      200
// @Failure      400  "Токен не передан или невалиден"
// @Failure      500  "Ошибка сервера"
// @Router       /api/sign-out [post]
func (h *Handler) Logout(gCtx *gin.Context) {
	// получаем заголовок
	jwtStr := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса, токен передан неверно
		gCtx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// отрезаем префикс
	jwtStr = jwtStr[len(jwtPrefix):]

	token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.JWT.Token), nil
	})
	if err != nil || !token.Valid {
		logrus.Error(err)
		gCtx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// в блеклисте токен нужен только до своего истечения
	claims := token.Claims.(*ds.JWTClaims)
	ttl := time.Until(time.Unix(claims.ExpiresAt, 0))

	// сохраняем в блеклист redis
	if err := h.Redis.WriteJWTToBlacklist(gCtx.Request.Context(), jwtStr, ttl); err != nil {
		logrus.Error(err)
		gCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	gCtx.Status(http.StatusOK)
}
