package handler

import (
	"fmt"
	"gofi/lib"
	"gofi/lib/constant"
	"gofi/middleware"
	"gofi/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/masb0ymas/go-utils/pkg"
)

type UploadHandler struct {
	service *service.UploadService
}

func NewUploadHandler(service *service.UploadService) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}

func (h *UploadHandler) RegisterRoutes(route fiber.Router) {
	// only admin can access
	adminOnly := []string{constant.ID_SUPER_ADMIN, constant.ID_ADMIN}

	new_route := route.Group("/upload")
	new_route.Get("/", h.GetAllUploads)
	new_route.Get("/:id", h.GetUploadById)
	new_route.Post("/", h.CreateUpload)
	new_route.Post("/public-view", h.CreateUploadPublicView)
	new_route.Put("/:id", middleware.Authorization(), middleware.PermissionAccess(adminOnly), h.UpdateUpload)
	new_route.Put("/restore/:id", middleware.Authorization(), middleware.PermissionAccess(adminOnly), h.RestoreUpload)
	new_route.Delete("/soft-delete/:id", middleware.Authorization(), middleware.PermissionAccess(adminOnly), h.SoftDeleteUpload)
	new_route.Delete("/force-delete/:id", middleware.Authorization(), middleware.PermissionAccess(adminOnly), h.ForceDeleteUpload)
}

func (h *UploadHandler) GetAllUploads(c *fiber.Ctx) error {
	req := &lib.Pagination{
		Page:     pkg.StringToInt32(c.Query("page", "1")),
		PageSize: pkg.StringToInt32(c.Query("page_size", "10")),
		Filtered: lib.ParseFilterItems(c.Query("filtered", "[]")),
		Sorted:   lib.ParseSortItems(c.Query("sorted", "[]")),
	}

	records, total, err := h.service.GetAllUploads(req)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendGetResponse(c, "data has been received", records, lib.Paginate(req.Page, req.PageSize, total))
}

func (h *UploadHandler) GetUploadById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid upload id")
	}

	record, err := h.service.GetUploadById(id)
	if err != nil {
		return lib.SendNotFoundResponse(c, "upload not found")
	}

	return lib.SendSuccessResponse(c, "data has been received", record)
}

func (h *UploadHandler) CreateUpload(c *fiber.Ctx) error {
	var req service.UploadRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	// Parse the form file
	file, err := c.FormFile("file_upload")
	if err != nil {
		return lib.SendBadRequestResponse(c, "failed to parse form file")
	}

	// Ensure file is not nil before proceeding
	if file.Size == 0 {
		return lib.SendBadRequestResponse(c, "uploaded file is empty")
	}

	fileData, err := file.Open()
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, fmt.Errorf("failed to open file: %w", err))
	}
	defer fileData.Close()

	keyFile := fmt.Sprintf("uploads/%s", uuid.New().String())
	newFilename := fmt.Sprintf("%s-%s", keyFile, file.Filename)

	signedURL, expiresAt, err := lib.UploadFile(newFilename, &fileData)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, fmt.Errorf("failed to upload file: %w", err))
	}

	uploadReq := service.CreateUploadRequest{
		KeyFile:   keyFile,
		Filename:  newFilename,
		Mimetype:  file.Header.Get("Content-Type"),
		Size:      file.Size,
		SignedURL: signedURL,
		ExpiresAt: expiresAt,
	}

	record, err := h.service.CreateUpload(&uploadReq)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendCreatedResponse(c, "data has been added", record)
}

func (h *UploadHandler) CreateUploadPublicView(c *fiber.Ctx) error {
	var req service.UploadRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	dirPath := req.ObjectPath
	if req.ObjectPath == nil {
		defaultPath := "uploads"
		dirPath = &defaultPath
	}

	// Parse the form file
	file, err := c.FormFile("file_upload")
	if err != nil {
		return lib.SendBadRequestResponse(c, "failed to parse form file")
	}

	// Ensure file is not nil before proceeding
	if file.Size == 0 {
		return lib.SendBadRequestResponse(c, "uploaded file is empty")
	}

	fileData, err := file.Open()
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, fmt.Errorf("failed to open file: %w", err))
	}
	defer fileData.Close()

	keyFile := fmt.Sprintf("%s/%s", *dirPath, uuid.New().String())
	newFilename := fmt.Sprintf("%s-%s", keyFile, file.Filename)

	publicURL, expiresAt, err := lib.UploadFilePublicView(newFilename, &fileData)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, fmt.Errorf("failed to upload file: %w", err))
	}

	uploadReq := service.CreateUploadRequest{
		KeyFile:   keyFile,
		Filename:  newFilename,
		Mimetype:  file.Header.Get("Content-Type"),
		Size:      file.Size,
		SignedURL: publicURL,
		ExpiresAt: expiresAt,
	}

	record, err := h.service.CreateUpload(&uploadReq)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendCreatedResponse(c, "data has been added", record)
}

func (h *UploadHandler) UpdateUpload(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid upload id")
	}

	var req service.UploadRequest
	if err := c.BodyParser(&req); err != nil {
		return lib.SendBadRequestResponse(c, "invalid request body")
	}

	// Parse the form file
	file, err := c.FormFile("file_upload")
	if err != nil {
		return lib.SendBadRequestResponse(c, "failed to parse form file")
	}

	// Ensure file is not nil before proceeding
	if file.Size == 0 {
		return lib.SendBadRequestResponse(c, "uploaded file is empty")
	}

	fileData, err := file.Open()
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, fmt.Errorf("failed to open file: %w", err))
	}
	defer fileData.Close()

	keyFile := fmt.Sprintf("uploads/%s", uuid.New().String())
	newFilename := fmt.Sprintf("%s-%s", keyFile, file.Filename)

	signedURL, expiresAt, err := lib.UploadFile(newFilename, &fileData)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, fmt.Errorf("failed to upload file: %w", err))
	}

	uploadReq := service.UpdateUploadRequest{
		KeyFile:   keyFile,
		Filename:  newFilename,
		Mimetype:  file.Header.Get("Content-Type"),
		Size:      file.Size,
		SignedURL: signedURL,
		ExpiresAt: expiresAt,
	}

	record, err := h.service.UpdateUpload(id, &uploadReq)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been updated", record)
}

func (h *UploadHandler) RestoreUpload(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid upload id")
	}

	err = h.service.RestoreUpload(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been restored", nil)
}

func (h *UploadHandler) SoftDeleteUpload(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid upload id")
	}

	err = h.service.SoftDeleteUpload(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been deleted", nil)
}

func (h *UploadHandler) ForceDeleteUpload(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return lib.SendBadRequestResponse(c, "invalid upload id")
	}

	err = h.service.ForceDeleteUpload(id)
	if err != nil {
		return lib.SendInternalServerErrorResponse(c, err)
	}

	return lib.SendSuccessResponse(c, "data has been permanently deleted from the system", nil)
}
