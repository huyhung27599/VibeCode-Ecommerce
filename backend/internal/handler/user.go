package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibecode/ecommerce/backend/internal/domain"
	"github.com/vibecode/ecommerce/backend/internal/dto"
	"github.com/vibecode/ecommerce/backend/internal/service"
	"github.com/vibecode/ecommerce/backend/pkg/response"
)

type User struct {
	svc service.UserService
}

func NewUser(svc service.UserService) *User {
	return &User{svc: svc}
}

func (h *User) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	u, err := h.svc.Create(c.Request.Context(), service.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     domain.Role(req.Role),
	})
	if err != nil {
		h.mapError(c, err)
		return
	}
	response.Created(c, dto.NewUserResponse(u))
}

func (h *User) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	u, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.mapError(c, err)
		return
	}
	response.OK(c, dto.NewUserResponse(u))
}

func (h *User) List(c *gin.Context) {
	var q dto.ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "invalid query params", err.Error())
		return
	}

	users, total, err := h.svc.List(c.Request.Context(), q.Page, q.PageSize)
	if err != nil {
		response.Internal(c, "failed to list users")
		return
	}

	totalPages := int(total) / q.PageSize
	if int(total)%q.PageSize != 0 {
		totalPages++
	}
	response.Paginated(c, dto.NewUserResponses(users), response.Meta{
		Page:       q.Page,
		PageSize:   q.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *User) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	in := service.UpdateUserInput{FullName: req.FullName, IsActive: req.IsActive}
	if req.Role != nil {
		r := domain.Role(*req.Role)
		in.Role = &r
	}

	u, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		h.mapError(c, err)
		return
	}
	response.OK(c, dto.NewUserResponse(u))
}

func (h *User) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.mapError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func (h *User) Me(c *gin.Context) {
	raw := c.GetString("user_id")
	id, err := uuid.Parse(raw)
	if err != nil {
		response.Unauthorized(c, "invalid session")
		return
	}
	u, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.mapError(c, err)
		return
	}
	response.OK(c, dto.NewUserResponse(u))
}

func (h *User) mapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, service.ErrEmailAlreadyInUse):
		response.Fail(c, 409, "CONFLICT", err.Error())
	case errors.Is(err, service.ErrInvalidCredential):
		response.Unauthorized(c, err.Error())
	default:
		response.Internal(c, "internal error")
	}
}
