package net

import "errors"

var (
	ErrTooShort      = errors.New("too short")
	ErrTooLong       = errors.New("too long")
	ErrWrongChecksum = errors.New("checksum didn't match")
)
