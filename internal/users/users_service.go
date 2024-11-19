package users

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	signingTokenKey = "joASdeDS3i#kjmFDSk3i303904lXSDds"
	tokenExpire     = 12 * time.Hour
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) CreateUser(user User) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service create user: %w", err)
	}
	user.Password = string(passwordHash)
	return us.repo.CreateUser(user)
}

type tokenClaims struct {
	ID       int
	Username string
	jwt.StandardClaims
}

func (us *UserService) Login(email, password string) (*UserLogin, error) {
	user, err := us.repo.Login(email)
	if err != nil {
		return nil, fmt.Errorf("service login: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("service login: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		user.ID,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	accessToken, err := token.SignedString([]byte(signingTokenKey))
	if err != nil {
		return nil, fmt.Errorf("service login: %w", err)
	}

	return &UserLogin{accessToken, user.ID, user.Email, user.Username}, nil
}

func (us *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(signingTokenKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("parse token: %w", err)
	}

	return claims.ID, nil
}
