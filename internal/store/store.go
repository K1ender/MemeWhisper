package store

import "errors"

var (
	ErrFailedToStartTransaction  = errors.New("failed to start transaction")
	ErrFailedToCommitTransaction = errors.New("failed to commit transaction")
	ErrFailedToHashPassword      = errors.New("failed to hash password")
	ErrUserAlreadyExists         = errors.New("user already exists")
)

var (
	ErrUserDoesntExist        = errors.New("user doesn't exist")
	ErrFailedToScanUser       = errors.New("failed to scan user")
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToUpdateUsername = errors.New("failed to update username")
	ErrFailedToUpdatePassword = errors.New("failed to update password")
	ErrFailedToGetUserID      = errors.New("failed to get user id")
)

var (
	ErrFailedToGenerateRandomBytes = errors.New("failed to generate random bytes")
	ErrFailedToCreateSession       = errors.New("failed to create session")
	ErrFailedToScanSession         = errors.New("failed to scan session")
	ErrSessionExpired              = errors.New("session expired")
	ErrFailedToDeleteSession       = errors.New("failed to delete session")
)
