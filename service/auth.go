package service

import (
	"errors"
	"fmt"
	"gofi/config"
	"gofi/database/model"
	"gofi/lib"
	"gofi/lib/constant"
	"gofi/repository"
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
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type AuthSignInRequest struct {
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required"`
	IpAddress string `json:"ip_address" form:"ip_address" validate:"required"`
	UserAgent string `json:"user_agent" form:"user_agent" validate:"required"`
}

func (s *AuthService) SignUp(req AuthSignUpRequest) (uuid.UUID, *string, error) {
	user := &model.User{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: &req.Password,
		RoleID:   uuid.MustParse(constant.ID_USER),
	}

	token, err := s.repo.SignUp(user)
	return user.ID, token, err
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
	token, err := lib.ExtractToken(c)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(c.Locals("uid").(string))
	if err != nil {
		return nil, fmt.Errorf("invalid UID in token: %v", err)
	}

	return s.repo.VerifySession(uid, token)
}

func (s *AuthService) VerifyToken(uid uuid.UUID, token string) error {
	return s.repo.VerifyToken(uid, token)
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

	uid, err := uuid.Parse(c.Locals("uid").(string))
	if err != nil {
		return fmt.Errorf("invalid UID in token: %v", err)
	}

	if claims.Exp < time.Now().Unix() {
		return errors.New("token is invalid")
	}

	return s.repo.SignOut(uid, token)
}
