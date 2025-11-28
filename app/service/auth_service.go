package service

import (
	"ekrp/app/models"
	"ekrp/app/repository"
	"ekrp/utils"
	"fmt"
)

type AuthService struct {
	UserRepo       *repository.UserRepository
	PermissionRepo *repository.PermissionRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepo:       userRepo,
		PermissionRepo: repository.NewPermissionRepository(),
	}
}

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// ambil user dari email
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// cek password
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 🔥 ambil permissions berdasarkan role
	perms, err := s.PermissionRepo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		return nil, err
	}

	// 🔥 generate token lengkap dgn permissions
	token, err := utils.GenerateTokenWithPermissions(user.ID, user.RoleID, perms)
	if err != nil {
		return nil, err
	}

	// response sesuai modul
	res := &models.LoginResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		FullName:   user.FullName,
		RoleID:     user.RoleID,
		Token:      token,
		Permissions: perms,
	}

	return res, nil
}
