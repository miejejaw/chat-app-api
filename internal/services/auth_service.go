package services

import (
	"chat-app-api/internal/repositories"
	"chat-app-api/internal/utils"
	"strconv"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID              uint   `json:"id"`
		Username        string `json:"username"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	} `json:"user"`
}

type AuthService interface {
	Login(username, password string) (LoginResponse, error)
	RenewAccessToken(refreshToken string) (string, error)
}

type AuthServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &AuthServiceImpl{userRepo: userRepo}
}

func (s *AuthServiceImpl) Login(username, password string) (LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return LoginResponse{}, err
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return LoginResponse{}, err
	}

	userClaims := utils.UserClaim{
		UserID:    strconv.FormatUint(uint64(user.ID), 10),
		Username:  user.Username,
		Email:     user.Email,
		TokenType: "access",
	}

	accessToken, refreshToken, err := utils.GenerateAccessAndRefreshTokens(userClaims)
	if err != nil {
		return LoginResponse{}, err
	}

	var response LoginResponse
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName
	response.User.ProfileImageUrl = user.ProfileImageUrl

	return response, nil
}

func (s *AuthServiceImpl) RenewAccessToken(refreshToken string) (string, error) {
	return utils.RenewAccessToken(refreshToken)
}
