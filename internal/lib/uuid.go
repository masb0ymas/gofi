package lib

import "github.com/google/uuid"

// ParseOptionalUUID parses s as a UUID, treating an empty string as absent so
// optional filters can be passed straight through.
func ParseOptionalUUID(s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
