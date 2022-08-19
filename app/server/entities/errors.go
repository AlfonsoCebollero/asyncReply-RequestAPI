package entities

import "errors"

var (
	ValidationError = errors.New("the received request is not valid")
	BadWorkflow     = errors.New("workflow does not exist")
)
