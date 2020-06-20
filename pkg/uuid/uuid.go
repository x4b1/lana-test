package uuid

import (
	"github.com/google/uuid"
)

func New() string {
	return uuid.Must(uuid.NewRandom()).String()
}
