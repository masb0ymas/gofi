package service

import (
	"errors"
	"goarif-api/config"
	"goarif-api/database/model"
	"goarif-api/lib"
	"goarif-api/repository"
	"time"

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

type CreateSessionRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type UpdateSessionRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (s *SessionService) GetAllSessions() ([]model.Session, error) {
	return s.repo.FindAll()
}

func (s *SessionService) GetSessionById(id uuid.UUID) (*model.Session, error) {
	session, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}
	return session, nil
}

func (s *SessionService) CreateSession(req *CreateSessionRequest) (*model.Session, error) {
	token, _, err := lib.GenerateToken(&lib.Payload{
		UID:       req.UserID,
		SecretKey: config.Env("JWT_SECRET_KEY", "secret"),
		ExpiresAt: config.Env("JWT_EXPIRES_IN", "30"), // days
	})
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID: req.UserID,
		Token:  token,
	}

	err = s.repo.Create(session)
	if err != nil {
		return nil, err
	}

	return s.GetSessionById(session.ID)
}

func (s *SessionService) UpdateSession(id uuid.UUID, req *UpdateSessionRequest) (*model.Session, error) {
	token, _, err := lib.GenerateToken(&lib.Payload{
		UID:       req.UserID,
		SecretKey: config.Env("JWT_SECRET_KEY", "secret"),
		ExpiresAt: config.Env("JWT_EXPIRES_IN", "30"), // days
	})
	if err != nil {
		return nil, err
	}

	session, err := s.GetSessionById(id)
	if err != nil {
		return nil, err
	}

	session.UserID = req.UserID
	session.Token = token

	err = s.repo.Update(session)
	if err != nil {
		return nil, err
	}

	return s.GetSessionById(session.ID)
}

func (s *SessionService) RestoreSession(id uuid.UUID) error {
	session, err := s.GetSessionById(id)
	if err != nil {
		return err
	}

	return s.repo.Restore(session.ID)
}

func (s *SessionService) SoftDeleteSession(id uuid.UUID) error {
	session, err := s.GetSessionById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(session.ID)
}

func (s *SessionService) ForceDeleteSession(id uuid.UUID) error {
	session, err := s.GetSessionById(id)
	if err != nil {
		return err
	}

	return s.repo.ForceDelete(session.ID)
}
