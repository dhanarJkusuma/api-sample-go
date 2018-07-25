package models

import "errors"

var (
	INTERNAL_SERVER_ERROR = errors.New("Internal Server Error")
	NOT_FOUND_ERROR       = errors.New("Resource Not Found")
	CONFLIT_ERROR         = errors.New("Resource is already exist")
	INVALID_REQUEST_ERROR = errors.New("Invalid Request")
)
