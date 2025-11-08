package handlers

import "gofi/internal/app"

type Handlers struct {
	Health healthHandler
	Role   roleHandler
}

func New(app *app.Application) Handlers {
	return Handlers{
		Health: healthHandler{app: app},
		Role:   roleHandler{app: app},
	}
}
