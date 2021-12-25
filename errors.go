package aternos_api

import "errors"

var ServerAlreadyStartedError = errors.New("server already running")
var ServerAlreadyStoppedError = errors.New("server already stopped")
