package apdex

import (
	"errors"
)

var (
	ErrInvalidIndex  = errors.New("invalid index")
	ErrHostNotExists = errors.New("host does not exists")
	ErrRuleMismatch  = errors.New("rule mismatch")
)
