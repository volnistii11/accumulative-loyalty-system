package repository

import "errors"

var (
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInternalServer       = errors.New("internal server error")

	ErrLoginOrPasswordIsIncorrect = errors.New("login or password is incorrect")
	ErrUserAlreadyRegistered      = errors.New("user already registered")
	ErrUserNotAuthorized          = errors.New("user not authorized")

	ErrOrderNumberHasAlreadyBeenUploadedByThisUser    = errors.New("order number has already been uploaded by this user")
	ErrOrderNumberHasAlreadyBeenUploadedByAnotherUser = errors.New("order number has already been uploaded by another user")

	ErrNotEnoughFundsOnTheAccount = errors.New("not enough funds on the account")
	ErrInvalidOrderNumber         = errors.New("invalid order number")

	ErrTheOrderIsNotRegisteredInTheBillingSystem      = errors.New("the order is not registered in the billing system")
	ErrTheNumberOfRequestsToTheServiceHasBeenExceeded = errors.New("the number of requests to the service has been exceeded")
)
