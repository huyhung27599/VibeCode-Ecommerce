package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vibecode/ecommerce/backend/internal/domain"
	"github.com/vibecode/ecommerce/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrInvalidCredential = errors.New("invalid credentials")
)

type CreateUserInput struct {
	Email    string
	Password string
	FullName string
	Role     domain.Role
}

type UpdateUserInput struct {
	FullName *string
	Role     *domain.Role
	IsActive *bool
}

type UserService interface {
	Create(ctx context.Context, in CreateUserInput) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateUserInput) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	VerifyPassword(ctx context.Context, email, password string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, in CreateUserInput) (*domain.User, error) {
	if existing, err := s.repo.GetByEmail(ctx, in.Email); err == nil && existing != nil {
		return nil, ErrEmailAlreadyInUse
	} else if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := in.Role
	if role == "" {
		role = domain.RoleUser
	}

	u := &domain.User{
		Email:        in.Email,
		PasswordHash: string(hash),
		FullName:     in.FullName,
		Role:         role,
		IsActive:     true,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, ErrEmailAlreadyInUse
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.List(ctx, (page-1)*pageSize, pageSize)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, in UpdateUserInput) (*domain.User, error) {
	u, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.FullName != nil {
		u.FullName = *in.FullName
	}
	if in.Role != nil {
		u.Role = *in.Role
	}
	if in.IsActive != nil {
		u.IsActive = *in.IsActive
	}
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *userService) VerifyPassword(ctx context.Context, email, password string) (*domain.User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredential
		}
		return nil, err
	}
	if !u.IsActive {
		return nil, ErrInvalidCredential
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredential
	}
	return u, nil
}
