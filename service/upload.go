package service

import (
	"database/sql"
	"errors"
	"gofi/database/model"
	"gofi/lib"
	"gofi/repository"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type UploadService struct {
	repo *repository.UploadRepository
}

func NewUploadService(repo *repository.UploadRepository) *UploadService {
	return &UploadService{
		repo: repo,
	}
}

type CreateUploadRequest struct {
	KeyFile   string    `json:"key_file" form:"key_file" validate:"required"`
	Filename  string    `json:"filename" form:"filename" validate:"required"`
	Mimetype  string    `json:"mimetype" form:"mimetype" validate:"required"`
	Size      int64     `json:"size" form:"size" validate:"required"`
	SignedURL string    `json:"signed_url" form:"signed_url" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" form:"expires_at" validate:"required"`
}

type UpdateUploadRequest struct {
	KeyFile   string    `json:"key_file" form:"key_file"`
	Filename  string    `json:"filename" form:"filename" validate:"required"`
	Mimetype  string    `json:"mimetype" form:"mimetype" validate:"required"`
	Size      int64     `json:"size" form:"size" validate:"required"`
	SignedURL string    `json:"signed_url" form:"signed_url" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" form:"expires_at" validate:"required"`
}

type UploadRequest struct {
	ObjectPath *string               `json:"object_path" form:"object_path"`
	FileUpload *multipart.FileHeader `json:"file_upload" form:"file_upload" validate:"required"`
}

func (s *UploadService) GetAllUploads(req *lib.Pagination) ([]model.Upload, int64, error) {
	return s.repo.FindAllWithPagination(req)
}

func (s *UploadService) GetUploadById(id uuid.UUID) (*model.Upload, error) {
	record, err := s.repo.FindById(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("role not found")
	}
	return record, nil
}

func (s *UploadService) GetUploadRecordById(id uuid.UUID) (*model.Upload, error) {
	record, err := s.repo.FindRecordById(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if record == nil {
		return nil, errors.New("role not found")
	}
	return record, nil
}

func (s *UploadService) CreateUpload(req *CreateUploadRequest) (*model.Upload, error) {
	values := &model.Upload{
		BaseModel: model.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		KeyFile:   req.KeyFile,
		Filename:  req.Filename,
		Mimetype:  req.Mimetype,
		Size:      req.Size,
		SignedURL: req.SignedURL,
		ExpiresAt: req.ExpiresAt,
	}

	err := s.repo.Create(values)
	if err != nil {
		return nil, err
	}

	return s.GetUploadById(values.ID)
}

func (s *UploadService) UpdateUpload(id uuid.UUID, req *UpdateUploadRequest) (*model.Upload, error) {
	values, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	values.KeyFile = req.KeyFile
	values.Filename = req.Filename
	values.Mimetype = req.Mimetype
	values.Size = req.Size
	values.SignedURL = req.SignedURL
	values.ExpiresAt = req.ExpiresAt

	err = s.repo.Update(values)
	if err != nil {
		return nil, err
	}

	return s.GetUploadById(values.ID)
}

func (s *UploadService) RestoreUpload(id uuid.UUID) error {
	record, err := s.GetUploadRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.Restore(record.ID)
}

func (s *UploadService) SoftDeleteUpload(id uuid.UUID) error {
	record, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(record.ID)
}

func (s *UploadService) ForceDeleteUpload(id uuid.UUID) error {
	record, err := s.GetUploadRecordById(id)
	if err != nil {
		return err
	}

	return s.repo.ForceDelete(record.ID)
}
