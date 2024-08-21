package domain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	validate        *validator.Validate
	ErrUserNotFound = errors.New("user with such credentials not found")
)

type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"loggin"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type UserIDKey string

const UserIDKeyForContext UserIDKey = "userID"

func init() {
	validate = validator.New()
}

type SighUpAndInInput struct {
	Login    string `json:"login" validate:"required,gte=2"`
	Password string `json:"password" validate:"required,gte=4"`
}

func (i *SighUpAndInInput) Validate() error {
	return validate.Struct(i)
}
