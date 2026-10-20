package ds

import (
	"RenderTimeEstimator/internal/app/role"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID   uint      `json:"user_id"`
	UserUUID uuid.UUID `json:"user_uuid"`
	Role     role.Role `json:"role"`
}