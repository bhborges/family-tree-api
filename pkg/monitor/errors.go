// Package monitor APM package
package monitor

import "errors"

var (
	// ErrEnvConfig is returned if some error occurs setting up the environent vars.
	ErrEnvConfig = errors.New("monitor: unable to setup environment variables")
	// ErrProvideApp is returned when unable to provide a monitor app.
	ErrProvideApp = errors.New("monitor: unable to provide app")
)
