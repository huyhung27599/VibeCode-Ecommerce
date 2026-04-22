package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/vibecode/ecommerce/backend/internal/domain"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	FullName string `json:"full_name" binding:"required,min=1,max=255"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name" binding:"omitempty,min=1,max=255"`
	Role     *string `json:"role" binding:"omitempty,oneof=user admin"`
	IsActive *bool   `json:"is_active"`
}

type ListUsersQuery struct {
	Page     int `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func NewUserResponses(users []domain.User) []UserResponse {
	out := make([]UserResponse, len(users))
	for i := range users {
		out[i] = NewUserResponse(&users[i])
	}
	return out
}
