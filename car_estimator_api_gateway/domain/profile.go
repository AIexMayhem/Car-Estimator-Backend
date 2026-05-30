package domain

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	FullName string	`json:"fullName"`
	Email string	`json:"email"`
	Phone string	`json:"phone"`	
	Password string	`json:"password"`
	BirthDate string	`json:"birthDate"`
}

type RegisterResponse struct {
	UserId uuid.UUID `json:"userId"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Access string	`json:"accessToken"`
	Refresh string	`json:"refreshToken"`
}

type LoginResponse struct {
	UserId uuid.UUID `json:"userId"`
	TokenResponse
}

type UserResponse struct {
	Id uuid.UUID `json:"userId"`
	FullName string `json:"fullName"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	BirthDate time.Time `json:"birthDate"`
	RegisterDate time.Time `json:"registerDate"`
}