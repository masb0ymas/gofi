package service

import (
	"database/sql"
	"errors"
	"gofi/config"
	"gofi/database/model"
	"gofi/lib"
	"gofi/repository"
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
	UserID uuid.UUID `json:"user_id" form:"user_id" validate:"required"`
}

type UpdateSessionRequest struct {
	UserID uuid.UUID `json:"user_id" form:"user_id" validate:"required"`
}

func (s *SessionService) GetAllSessions(req *lib.Pagination) ([]model.Session, int64, error) {
	return s.repo.FindAllWithPagination(req)
}

func (s *SessionService) GetSessionById(id uuid.UUID) (*model.Session, error) {
	record, err := s.repo.FindById(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("session not found")
	}
	return record, nil
}

func (s *SessionService) GetSessionRecordById(id uuid.UUID) (*model.Session, error) {
	record, err := s.repo.FindRecordById(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("session not found")
	}
	return record, nil
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

	values := &model.Session{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID: req.UserID,
		Token:  token,
	}

	err = s.repo.Create(values)
	if err != nil {
		return nil, err
	}

	return s.GetSessionById(values.ID)
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

	values, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	values.UserID = req.UserID
	values.Token = token

	err = s.repo.Update(values)
	if err != nil {
		return nil, err
	}

	return s.GetSessionById(values.ID)
}

func (s *SessionService) RestoreSession(id uuid.UUID) error {
	record, err := s.GetSessionRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.Restore(record.ID)
}

func (s *SessionService) SoftDeleteSession(id uuid.UUID) error {
	record, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(record.ID)
}

func (s *SessionService) ForceDeleteSession(id uuid.UUID) error {
	record, err := s.GetSessionRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.ForceDelete(record.ID)
}
