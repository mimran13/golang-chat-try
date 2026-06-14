package dto

type LoginRequest struct {
    Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginResponse struct {
      Token string       `json:"token"`
      User  UserResponse `json:"user"`
  }