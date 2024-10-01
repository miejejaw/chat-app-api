package services

import (
	"chat-app-api/internal/repositories"
	"chat-app-api/internal/utils"
	"strconv"
)

type AuthService interface {
	Login(username, password string) (string, string, error)
	RenewAccessToken(refreshToken string) (string, error)
}

type AuthServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &AuthServiceImpl{userRepo: userRepo}
}

func (s *AuthServiceImpl) Login(username, password string) (string, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", "", err
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return "", "", err
	}

	userClaims := utils.UserClaim{
		UserID:    strconv.FormatUint(uint64(user.ID), 10),
		Username:  user.Username,
		Email:     user.Email,
		TokenType: "access",
	}

	accessToken, refreshToken, err := utils.GenerateAccessAndRefreshTokens(userClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthServiceImpl) RenewAccessToken(refreshToken string) (string, error) {
	return utils.RenewAccessToken(refreshToken)
}
