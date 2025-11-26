package service

import (
	"ekrp/app/models"
	"ekrp/app/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		Repo: repository.NewUserRepository(),
	}
}

func (s *UserService) CreateUser(username, email, password, fullName, role_id string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		RoleID:       role_id,
		IsActive:     true,
	}

	return s.Repo.CreateUser(user)
}
