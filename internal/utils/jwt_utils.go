package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var (
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	accessTokenExp     time.Duration
	refreshTokenExp    time.Duration
)

func init() {
	//Check the environment (development or production)
	env := os.Getenv("GIN_MODE")
	if env == "" {
		env = gin.DebugMode // default to debug if GIN_MODE is not set
	}

	// If we're in debug mode (development), try to load .env file
	if env == gin.DebugMode {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Read environment variables
	accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

	accessTokenExpiration := os.Getenv("ACCESS_TOKEN_EXPIRATION")
	refreshTokenExpiration := os.Getenv("REFRESH_TOKEN_EXPIRATION")

	var err error
	accessTokenExp, err = time.ParseDuration(accessTokenExpiration)
	if err != nil {
		log.Fatalf("Error parsing ACCESS_TOKEN_EXPIRATION: %v", err)
	}

	refreshTokenExp, err = time.ParseDuration(refreshTokenExpiration)
	if err != nil {
		log.Fatalf("Error parsing REFRESH_TOKEN_EXPIRATION: %v", err)
	}
}

type UserClaim struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
}

type Claims struct {
	UserClaim
	jwt.StandardClaims
}

func generateToken(userClaims UserClaim, secret []byte, exp time.Duration) (string, error) {
	expirationTime := time.Now().Add(exp)
	claims := &Claims{
		UserClaim: userClaims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseToken(tokenString string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func ParseAccessToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, accessTokenSecret)
}

func RenewAccessToken(refreshTokenString string) (string, error) {
	// Parse and validate the refresh token
	claims, err := parseToken(refreshTokenString, refreshTokenSecret)

	if err != nil {
		return "", err
	}

	// Ensure the refresh token is valid and not expired
	if claims.ExpiresAt < time.Now().Unix() {
		return "", errors.New("refresh token has expired")
	}

	// Generate a new access token
	userClaims := UserClaim{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Email:     claims.Email,
		TokenType: "access",
	}
	return generateToken(userClaims, accessTokenSecret, accessTokenExp)
}

func GenerateAccessAndRefreshTokens(userClaims UserClaim) (string, string, error) {
	accessToken, err := generateToken(userClaims, accessTokenSecret, accessTokenExp)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateToken(userClaims, refreshTokenSecret, refreshTokenExp)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
