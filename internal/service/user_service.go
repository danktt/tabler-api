package service

import (
	"context"
	"fmt"

	"tabler-api/internal/model"
	"tabler-api/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Create new user
	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*model.UserResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var responses []*model.UserResponse
	for _, user := range users {
		response := user.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if email is being changed and if it already exists
	if user.Email != req.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser != nil {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
	}

	// Update user
	user.Name = req.Name
	user.Email = req.Email

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
} 