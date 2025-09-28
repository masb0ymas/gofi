package service

import (
	"database/sql"
	"errors"
	"gofi/src/app/database/model"
	"gofi/src/app/repository"
	"gofi/src/lib"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type CreateUserRequest struct {
	Fullname string    `json:"fullname" form:"fullname" validate:"required"`
	Email    string    `json:"email" form:"email" validate:"required"`
	Phone    *string   `json:"phone" form:"phone"`
	IsActive bool      `json:"is_active" form:"is_active" validate:"required"`
	RoleID   uuid.UUID `json:"role_id" form:"role_id" validate:"required"`
}

type UpdateUserRequest struct {
	Fullname string    `json:"fullname" form:"fullname"`
	Email    string    `json:"email" form:"email"`
	Phone    *string   `json:"phone" form:"phone"`
	IsActive bool      `json:"is_active" form:"is_active" validate:"required"`
	RoleID   uuid.UUID `json:"role_id" form:"role_id" validate:"required"`
}

func (s *UserService) GetAllUsers(req *lib.Pagination) ([]model.User, int64, error) {
	return s.repo.FindAllWithPagination(req)
}

func (s *UserService) GetUserById(id uuid.UUID) (*model.User, error) {
	record, err := s.repo.FindById(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("user not found")
	}
	return record, nil
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	record, err := s.repo.FindByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("user not found")
	}
	return record, nil
}

func (s *UserService) GetUserRecordById(id uuid.UUID) (*model.User, error) {
	record, err := s.repo.FindRecordById(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("user not found")
	}
	return record, nil
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*model.User, error) {
	values := &model.User{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Fullname: req.Fullname,
		Email:    req.Email,
		Phone:    req.Phone,
		IsActive: req.IsActive,
		RoleID:   req.RoleID,
	}

	err := s.repo.Create(values)
	if err != nil {
		return nil, err
	}

	return s.GetUserById(values.ID)
}

func (s *UserService) UpdateUser(id uuid.UUID, req *UpdateUserRequest) (*model.User, error) {
	values, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	values.Fullname = req.Fullname
	values.Email = req.Email
	values.Phone = req.Phone
	values.IsActive = req.IsActive
	values.RoleID = req.RoleID

	err = s.repo.Update(values)
	if err != nil {
		return nil, err
	}

	return s.GetUserById(values.ID)
}

func (s *UserService) RestoreUser(id uuid.UUID) error {
	record, err := s.GetUserRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.Restore(record.ID)
}

func (s *UserService) SoftDeleteUser(id uuid.UUID) error {
	record, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(record.ID)
}

func (s *UserService) ForceDeleteUser(id uuid.UUID) error {
	record, err := s.GetUserRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.ForceDelete(record.ID)
}
