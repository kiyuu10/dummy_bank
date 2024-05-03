package api

import (
	"time"

	"github.com/google/uuid"
)

type (
	CreateUserReq struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	UserResponse struct {
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	LoginReq struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
	}

	LoginRes struct {
		SessionID             uuid.UUID    `json:"session_id"`
		AccessToken           string       `json:"access_token"`
		AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at"`
		RefreshToken          string       `json:"refresh_token"`
		RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
		User                  UserResponse `json:"user"`
	}
)
