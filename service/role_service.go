package service

import (
	"errors"
	"goarif-api/database/model"
	"goarif-api/repository"
	"time"

	"github.com/google/uuid"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	return s.repo.FindAll()
}

func (s *RoleService) GetRoleById(id uuid.UUID) (*model.Role, error) {
	role, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (s *RoleService) CreateRole(req *CreateRoleRequest) (*model.Role, error) {
	role := &model.Role{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: req.Name,
	}

	err := s.repo.Create(role)
	if err != nil {
		return nil, err
	}

	return s.GetRoleById(role.ID)
}

func (s *RoleService) UpdateRole(id uuid.UUID, req *UpdateRoleRequest) (*model.Role, error) {
	role, err := s.GetRoleById(id)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name

	err = s.repo.Update(role)
	if err != nil {
		return nil, err
	}

	return s.GetRoleById(role.ID)
}

func (s *RoleService) RestoreRole(id uuid.UUID) error {
	role, err := s.GetRoleById(id)
	if err != nil {
		return err
	}

	return s.repo.Restore(role.ID)
}

func (s *RoleService) SoftDeleteRole(id uuid.UUID) error {
	role, err := s.GetRoleById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(role.ID)
}

func (s *RoleService) ForceDeleteRole(id uuid.UUID) error {
	role, err := s.GetRoleById(id)
	if err != nil {
		return err
	}

	return s.repo.ForceDelete(role.ID)
}
