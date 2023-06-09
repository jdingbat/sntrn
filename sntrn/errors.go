package sntrn

import "errors"

var (
	ErrFailedLogin = errors.New("failed login")

	ErrFailedLogout = errors.New("failed logout")

	ErrNotLoggedIn = errors.New("client not logged in")
)
