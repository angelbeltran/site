package model

import (
	"time"
)

type (
	User struct {
		ID UserID
		Name string
		CreatedAt time.Time
	}

	UserID int

	UserSession struct {
		ID UserSessionID
		UserID UserID
		ExpiresAt time.Time
	}

	UserSessionID int
)
