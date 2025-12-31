package common

import "github.com/google/uuid"

func ParseUUID(str string) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, BadRequest("Invalid UUID!")
	}
	return id, nil
}
