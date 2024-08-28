package service

import (
	"context"
	"gofi/database/entity"
	"gofi/database/repository"

	"github.com/google/uuid"
)

type SessionService struct {
	repo *repository.SessionRepository
}

func NewSessionService(repo *repository.SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, value *entity.Session) (*entity.Session, error) {
	return s.repo.CreateSession(ctx, value)
}

func (s *SessionService) GetSession(ctx context.Context, id uuid.UUID, token string) (*entity.Session, error) {
	return s.repo.GetSession(ctx, id, token)
}

func (s *SessionService) ListSessions(ctx context.Context) ([]entity.Session, error) {
	return s.repo.ListSessions(ctx)
}

func (s *SessionService) UpdateSession(ctx context.Context, value *entity.Session) (*entity.Session, error) {
	return s.repo.UpdateSession(ctx, value)
}

func (s *SessionService) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSession(ctx, id)
}
