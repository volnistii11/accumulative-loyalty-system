package customErrors

import "github.com/pkg/errors"

var ErrHTTPStatusTooManyRequests = errors.New("too many requests")
var ErrHTTPStatusNoContent = errors.New("no content")
