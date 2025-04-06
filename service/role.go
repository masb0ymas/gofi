package service

import (
	"errors"
	"gofi/database/model"
	"gofi/lib"
	"gofi/repository"
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
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

func (s *RoleService) GetAllRoles(req *lib.Pagination) ([]model.Role, int64, error) {
	return s.repo.FindAllWithPagination(req)
}

func (s *RoleService) GetRoleById(id uuid.UUID) (*model.Role, error) {
	record, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("role not found")
	}
	return record, nil
}

func (s *RoleService) GetRoleRecordById(id uuid.UUID) (*model.Role, error) {
	record, err := s.repo.FindRecordById(id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("role not found")
	}
	return record, nil
}

func (s *RoleService) CreateRole(req *CreateRoleRequest) (*model.Role, error) {
	values := &model.Role{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: req.Name,
	}

	err := s.repo.Create(values)
	if err != nil {
		return nil, err
	}

	return s.GetRoleById(values.ID)
}

func (s *RoleService) UpdateRole(id uuid.UUID, req *UpdateRoleRequest) (*model.Role, error) {
	values, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	values.Name = req.Name

	err = s.repo.Update(values)
	if err != nil {
		return nil, err
	}

	return s.GetRoleById(values.ID)
}

func (s *RoleService) RestoreRole(id uuid.UUID) error {
	record, err := s.GetRoleRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.Restore(record.ID)
}

func (s *RoleService) SoftDeleteRole(id uuid.UUID) error {
	record, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(record.ID)
}

func (s *RoleService) ForceDeleteRole(id uuid.UUID) error {
	record, err := s.GetRoleRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.ForceDelete(record.ID)
}
