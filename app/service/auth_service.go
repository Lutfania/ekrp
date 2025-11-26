package service

import (
	"ekrp/app/models"
	"ekrp/app/repository"
	"ekrp/utils"
	"fmt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: repo}
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		// raw DB "no rows" -> ubah pesan (biar tidak bocor)
		return nil, fmt.Errorf("invalid credentials")
	}

	// compare using utils
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// generate jwt token
	token, err := utils.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return nil, err
	}

	res := &models.LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		RoleID:   user.RoleID,
		Token:    token,
	}
	return res, nil
}
