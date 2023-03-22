package http

import "errors"

var (
	// ErrEnvConfig is returned if some error occurs setting up the environent vars.
	ErrEnvConfig = errors.New("http: unable to setup environment variables")
	// ErrHTTPCloseConn is returned if unable to close http server connection.
	ErrHTTPCloseConn = errors.New("unable to close http server connection")
	// ErrTCPListening is returned if unable to announce to local network.
	ErrTCPListening = errors.New("unable to announce TCP connection")
)
