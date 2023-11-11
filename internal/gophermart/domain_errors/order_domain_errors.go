package domain_errors

import "errors"

var ErrOrderNumberNotValid = errors.New("order number not valid")
var ErrOrderLoadedByAnotherUser = errors.New("order loaded by another user")
var ErrOrderAlreadyLoaded = errors.New("order already loaded")
