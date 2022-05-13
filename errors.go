package aternos_api

import "errors"

var (
	ServerAlreadyStartedError = errors.New("server already started")

	ServerAlreadyStoppedError = errors.New("server already stopped")

	// UnauthenticatedError indicates an invalid account was used to request the resource.
	UnauthenticatedError = errors.New("unauthenticated (invalid account)")

	// ForbiddenError indicates that the request was blocked by CloudFlare.
	ForbiddenError = errors.New("forbidden (blocked by CloudFlare)")
)
