package dtos

import (
	"github.com/gofiber/fiber/v3"
)

// Metadata mirrors the web UI metadata type: an open map with conventional
// keys (total, offset, limit, page, per_page, from, to).
type Metadata map[string]interface{}

// ListMeta builds the conventional offset/limit pagination metadata.
func ListMeta(total, offset, limit int) Metadata {
	return Metadata{"total": total, "offset": offset, "limit": limit}
}

// TotalMeta builds metadata carrying only the row count.
func TotalMeta(total int) Metadata {
	return Metadata{"total": total}
}

// Response renders the single-resource envelope {data, metadata, message?}.
func Response(c fiber.Ctx, status int, data any, metadata Metadata, message string) error {
	body := fiber.Map{"data": data, "metadata": metadata}
	if message != "" {
		body["message"] = message
	}
	return c.Status(status).JSON(body)
}

// OK renders a 200 single-resource envelope without a message.
func OK(c fiber.Ctx, data any) error {
	return Response(c, fiber.StatusOK, data, Metadata{}, "")
}

// Created renders a 201 single-resource envelope with a message.
func Created(c fiber.Ctx, data any, message string) error {
	return Response(c, fiber.StatusCreated, data, Metadata{}, message)
}

// List renders the collection envelope {data, metadata}.
func List(c fiber.Ctx, data any, meta Metadata) error {
	return Response(c, fiber.StatusOK, data, meta, "data has been retrieved")
}

// Get renders the single-resource envelope {data}.
func Get(c fiber.Ctx, data any) error {
	return Response(c, fiber.StatusOK, data, Metadata{}, "data has been retrieved")
}

// Message renders the bare {message} envelope (accepted for updates that return no payload).
func Message(c fiber.Ctx, status int, message string) error {
	return Response(c, status, nil, Metadata{}, message)
}
