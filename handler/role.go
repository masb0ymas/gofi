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

type roleHandler struct {
	ctx     context.Context
	service *service.RoleService
}

func NewRoleHandler(service *service.RoleService) *roleHandler {
	return &roleHandler{
		ctx:     context.Background(),
		service: service,
	}
}

func toStoreRole(r *entity.RoleReq) *entity.Role {
	return &entity.Role{
		Name: r.Name,
	}
}

func toRoleRes(r *entity.Role) entity.RoleRes {
	return entity.RoleRes{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
		Name:      r.Name,
	}
}

func pathRoleReq(role *entity.Role, r entity.RoleReq) {
	if r.Name != "" {
		role.Name = r.Name
	}

	role.UpdatedAt = toTimePtr(time.Now())
}

func (h *roleHandler) createRole(c *fiber.Ctx) error {
	r := new(entity.RoleReq)
	if code, message, errors := utils.ParseFormDataAndValidate(c, r); errors != nil {
		response := utils.FailureResponse(code, message, errors)
		return c.Status(int(code)).JSON(response)
	}

	record, err := h.service.CreateRole(h.ctx, toStoreRole(r))
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been added", record)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *roleHandler) getRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errFiber := fiber.NewError(http.StatusBadRequest)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	record, err := h.service.GetRole(h.ctx, id)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been received", record)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *roleHandler) listRoles(c *fiber.Ctx) error {
	roles, err := h.service.ListRoles(h.ctx)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	var res []entity.RoleRes
	for _, p := range roles {
		res = append(res, toRoleRes(&p))
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been received", res)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *roleHandler) updateRole(c *fiber.Ctx) error {
	// parsing uuid
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errFiber := fiber.NewError(http.StatusBadRequest)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	// form validation
	r := new(entity.RoleReq)
	if code, message, errors := utils.ParseFormDataAndValidate(c, r); errors != nil {
		response := utils.FailureResponse(code, message, errors)
		return c.Status(int(code)).JSON(response)
	}

	// get role by id
	role, err := h.service.GetRole(h.ctx, id)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	// path update
	pathRoleReq(role, *r)
	updated, err := h.service.UpdateRole(h.ctx, role)
	if err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been updated", updated)
	return c.Status(http.StatusOK).JSON(response)
}

func (h *roleHandler) deleteRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		errFiber := fiber.NewError(http.StatusBadRequest)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	if err := h.service.DeleteRole(h.ctx, id); err != nil {
		errFiber := fiber.NewError(http.StatusInternalServerError)
		response := utils.FailureResponse(int32(errFiber.Code), errFiber.Message, []string{err.Error()})
		return c.Status(errFiber.Code).JSON(response)
	}

	response := utils.SuccessResponse(http.StatusOK, "data has been deleted", nil)
	return c.Status(http.StatusOK).JSON(response)
}
