package err

import (
	"errors"
)

// ErrBadIP indicates ip is not valid
var ErrBadIP = errors.New("bad ip value")

// ErrBadPort indicates port is not valid or not authorized
var ErrBadPort = errors.New("bad port value")

// ErrBadFmt
var ErrBadFmt = errors.New("bad format")
