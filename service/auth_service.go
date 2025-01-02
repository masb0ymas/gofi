package service

import (
	"errors"
	"fmt"
	"goarif-api/config"
	"goarif-api/database/model"
	"goarif-api/lib"
	"goarif-api/lib/constant"
	"goarif-api/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

type AuthSignUpRequest struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthSignInRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	IpAddress string `json:"ip_address" validate:"required"`
	UserAgent string `json:"user_agent" validate:"required"`
}

func (s *AuthService) SignUp(req AuthSignUpRequest) error {
	user := &model.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: &req.Password,
		RoleID:   uuid.MustParse(constant.ID_USER),
	}
	return s.repo.SignUp(user)
}

func (s *AuthService) SignIn(req AuthSignInRequest) (*model.SignInResponse, error) {
	user := &model.User{
		Email:    req.Email,
		Password: &req.Password,
	}
	auth, err := s.repo.SignIn(user, req.IpAddress, req.UserAgent)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, errors.New("user not found")
	}

	return auth, nil
}

func (s *AuthService) VerifySession(c *fiber.Ctx) (*model.User, error) {
	secretKey := config.Env("JWT_SECRET_KEY", "secret")
	token, err := lib.ExtractToken(c)
	if err != nil {
		return nil, err
	}

	claims, err := lib.VerifyToken(c, secretKey)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(claims.UID)
	if err != nil {
		return nil, fmt.Errorf("invalid UID in token: %v", err)
	}

	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token is invalid")
	}

	return s.repo.VerifySession(uid, token)
}

func (s *AuthService) SignOut(c *fiber.Ctx) error {
	secretKey := config.Env("JWT_SECRET_KEY", "secret")
	token, err := lib.ExtractToken(c)
	if err != nil {
		return err
	}

	claims, err := lib.VerifyToken(c, secretKey)
	if err != nil {
		return err
	}

	uid, err := uuid.Parse(claims.UID)
	if err != nil {
		return fmt.Errorf("invalid UID in token: %v", err)
	}

	if claims.Exp < time.Now().Unix() {
		return errors.New("token is invalid")
	}

	return s.repo.SignOut(uid, token)
}
