package eosc

import "errors"

var (
	ErrorRegisterConflict = errors.New("conflict of register")
	ErrorConfigType       = errors.New("error config type")
)
