package dto

import "errors"

var (
	ErrBadRequest = errors.New("invalid input data") // 400
	ErrNotFound   = errors.New("not found")          // 404
)