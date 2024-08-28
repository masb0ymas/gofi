package service

import (
	"context"
	"gofi/database/entity"
	"gofi/database/repository"

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

func (s *RoleService) CreateRole(ctx context.Context, value *entity.Role) (*entity.Role, error) {
	return s.repo.CreateRole(ctx, value)
}

func (s *RoleService) GetRole(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	return s.repo.GetRole(ctx, id)
}

func (s *RoleService) ListRoles(ctx context.Context) ([]entity.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *RoleService) UpdateRole(ctx context.Context, value *entity.Role) (*entity.Role, error) {
	return s.repo.UpdateRole(ctx, value)
}

func (s *RoleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteRole(ctx, id)
}
