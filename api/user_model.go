package api

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
)
