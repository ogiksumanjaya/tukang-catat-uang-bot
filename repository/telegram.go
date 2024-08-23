package repository

import (
	"strings"
)

type TelegramRepository struct {
	usernameAllowed string
}

func NewTelegramRepository(usernameAllowed string) *TelegramRepository {
	return &TelegramRepository{
		usernameAllowed: usernameAllowed,
	}
}

func (t *TelegramRepository) IsAllowedUsername(username string) bool {
	allowedUsername := strings.Split(t.usernameAllowed, ",")
	for _, allowed := range allowedUsername {
		if username == allowed {
			return true
		}
	}
	return false
}
