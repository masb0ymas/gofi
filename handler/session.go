package handler

import (
	"context"
	"gofi/database/entity"
	"gofi/pkg/utils"
	"gofi/service"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type sessionHandler struct {
	ctx     context.Context
	service *service.SessionService
}

func NewSessionHandler(service *service.SessionService) *sessionHandler {
	return &sessionHandler{
		ctx:     context.Background(),
		service: service,
	}
}

func toStoreSession(s *entity.SessionReq) *entity.Session {
	return &entity.Session{
		UserID:    s.UserID,
		Token:     s.Token,
		ExpiredAt: s.ExpiredAt,
	}
}

func toSessionRes(s *entity.Session) entity.SessionRes {
	return entity.SessionRes{
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		UserID:    s.UserID,
		Token:     s.Token,
		ExpiredAt: s.ExpiredAt,
	}
}

func pathSessionReq(session *entity.Session, s entity.SessionReq) {
	if isValidUUID(s.UserID.String()) {
		session.UserID = s.UserID
	}

	if s.Token != "" {
		session.Token = s.Token
	}

	if isDateValue(s.ExpiredAt.String()) {
		session.ExpiredAt = s.ExpiredAt
	}

	session.UpdatedAt = toTimePtr(time.Now())
}

func (h *sessionHandler) createSession(c *fiber.Ctx) error {
	input := new(entity.SessionReq)
	if code, message, errors := utils.ParseFormDataAndValidate(c, input); errors != nil {
		response := utils.FailureResponse(code, message, errors)
		return c.Status(int(code)).JSON(response)
	}

	record, err := h.service.CreateSession(h.ctx, toStoreSession(input))
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been added", record)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *sessionHandler) getSession(c *fiber.Ctx) error {
	token := c.Params("token")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errFiber := fiber.NewError(http.StatusBadRequest)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	record, err := h.service.GetSession(h.ctx, id, token)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been received", record)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *sessionHandler) listSessions(c *fiber.Ctx) error {
	roles, err := h.service.ListSessions(h.ctx)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	var res []entity.SessionRes
	for _, p := range roles {
		res = append(res, toSessionRes(&p))
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been received", res)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *sessionHandler) updateSession(c *fiber.Ctx) error {
	// parsing uuid
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errFiber := fiber.NewError(http.StatusBadRequest)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	// form validation
	input := new(entity.SessionReq)
	if code, message, errors := utils.ParseFormDataAndValidate(c, input); errors != nil {
		response := utils.FailureResponse(code, message, errors)
		return c.Status(int(code)).JSON(response)
	}

	// get session by id
	session, err := h.service.GetSession(h.ctx, id, input.Token)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	// path update
	pathSessionReq(session, *input)
	updated, err := h.service.UpdateSession(h.ctx, session)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been updated", updated)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *sessionHandler) deleteSession(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errFiber := fiber.NewError(http.StatusBadRequest)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	if err := h.service.DeleteSession(h.ctx, id); err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been deleted", nil)
	return c.Status(http.StatusOK).JSON(response)
}
