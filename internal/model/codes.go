package model

import "errors"

var (
	ErrSuccess             = errors.New("success")               // 200 OK
	ErrCreated             = errors.New("created")               // 201 Created
	ErrNoContent           = errors.New("no content")            // 204 No Content
	ErrBadRequest          = errors.New("bad request")           // 400 Bad Request
	ErrNotFound            = errors.New("not found")             // 404 Not Found
	ErrMethodNotAllowed    = errors.New("method not allowed")    // 405 Method Not Allowed
	ErrInternalServerError = errors.New("internal server error") // 500 Internal Server Error
)
