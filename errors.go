package aternos_api

import "errors"

var ServerAlreadyStartedError = errors.New("server already started")
var ServerAlreadyStoppedError = errors.New("server already stopped")
var UnauthenticatedError = errors.New("unauthenticated")
