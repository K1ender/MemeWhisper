package store

import "errors"

var (
	ErrFailedToStartTransaction  = errors.New("failed to start transaction")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction")
	ErrFailedToHashPassword      = errors.New("failed to hash password")

	ErrUserDoesntExist        = errors.New("user doesn't exist")
	ErrFailedToScanUser       = errors.New("failed to scan user")
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToUpdateUsername = errors.New("failed to update username")
	ErrFailedToUpdatePassword = errors.New("failed to update password")
)
