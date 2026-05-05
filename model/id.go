package model

import "github.com/google/uuid"

func NewStringID() string {
	return uuid.NewString()
}
