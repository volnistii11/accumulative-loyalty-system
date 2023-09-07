package cerrors

import "github.com/pkg/errors"

var ErrHTTPStatusTooManyRequests = errors.New("too many requests")
var ErrHTTPStatusNoContent = errors.New("no content")
var ErrHTTPStatusUnauthorized = errors.New("user unauthorized")
var ErrDBOrderExistsAndBelongsToTheUser = errors.New("order exists and belongs to the user")
var ErrDBOrderExistsAndDoesNotBelongToTheUser = errors.New("order exists and not belongs to the user")
var ErrDBNotEnoughCoins = errors.New("not enough coins")
var ErrHTTPWrongRequestFormat = errors.New("wrong request format")
var ErrHTTPUserExists = errors.New("user exists")
var ErrInternalServer = errors.New("internal server error")
var ErrUserOrPasswordIncorrect = errors.New("user or pass is incorrect")
